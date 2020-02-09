package command

import (
	"fmt"
	"log"
	"reflect"

	"github.com/bengtrj/simple-sync/config"
	"github.com/bengtrj/simple-sync/sshclient"
)

//Sync runs the synchronization logic
func Sync(config *config.Sync) error {
	for _, server := range config.DesiredState.Servers {
		for _, desiredApp := range config.DesiredState.Apps {
			err := func() error {
				knownApp := findApp(config.KnownState, desiredApp.Name)
				if reflect.DeepEqual(knownApp, desiredApp) {
					fmt.Printf("Skipping app %s on server %s because it's already on the desired state.",
						knownApp.Name, server)
					return nil
				}
				prettyPrintSync(knownApp, desiredApp, server)
				address := fmt.Sprintf("%s:%s", server.IP, "22")
				client, err := sshclient.DialWithPasswd(address, config.User, config.Password)
				if err != nil {
					return err
				}
				defer client.Close()

				return synchronize(client, knownApp, desiredApp)
			}()
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func synchronize(client *sshclient.Client, knownApp, desiredApp config.App) error {
	err := stopServices(client, knownApp)
	if err != nil {
		return err
	}

	err = syncFiles(client, knownApp, desiredApp)
	if err != nil {
		return err
	}

	err = syncPackages(client, knownApp, desiredApp)
	if err != nil {
		return err
	}

	err = startServices(client, desiredApp)
	return err
}

func prettyPrintSync(knownApp, desiredApp config.App, server config.Server) {
	op := "Updating"
	if isNewInstall(knownApp) {
		op = "Installing"
	}
	fmt.Printf("\n****  %s app %s on server %s  ****\n", op, desiredApp.Name, server.IP)
}

func isNewInstall(app config.App) bool {
	return len(app.Packages) == 0
}

func findApp(knownState *config.State, name string) config.App {
	if knownState != nil {
		for _, app := range knownState.Apps {
			if app.Name == name {
				return app
			}
		}
	}
	// return a empty "known state"
	return config.App{}
}

// Installs all packages
func syncPackages(client *sshclient.Client, knownApp, desiredApp config.App) error {

	err := aptUpdate(client)
	if err != nil {
		return err
	}

	//First, add all known packages into a hashtable
	known := make(map[string]bool)
	for _, p := range knownApp.Packages {
		known[p.Name] = true
	}

	//Then:
	// if they are present on both, no-op
	// if they are only present on known, uninstall
	// if they are only present on desired, install
	for _, p := range desiredApp.Packages {

		if _, ok := known[p.Name]; ok {
			delete(known, p.Name)
			fmt.Printf("Skipping already installed package %s\n", p.Name)
		} else {
			fmt.Printf("Installing new package %s\n", p.Name)

			cmd := fmt.Sprintf("sudo apt-get install %s -y", p.Name)
			_, err := client.Cmd(cmd).SmartOutput()
			if err != nil {
				return err
			}
		}

	}

	//Now, all packages left on known should be uninstalled
	for name := range known {
		removePackage(client, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func removePackage(client *sshclient.Client, name string) error {
	fmt.Printf("Unistalling package %s\n", name)
	cmd := fmt.Sprintf("sudo apt-get remove %s -y", name)
	out, err := client.Cmd(cmd).SmartOutput()
	if err != nil {
		log.Print(string(out))
	}
	return err
}

// Syncs all files
func syncFiles(client *sshclient.Client, knownApp, desiredApp config.App) error {

	//First, add all known files into a hashtable
	known := make(map[string]bool)
	for _, p := range knownApp.Files {
		known[p.Path] = true
	}

	var err error

	//Then:
	// if they are only present on known, delete the file
	// otherwise copy overriding if exists
	for _, file := range desiredApp.Files {
		if _, ok := known[file.Path]; ok {
			delete(known, file.Path)
			fmt.Printf("Overriding file %s\n", file.Path)
			err = copyFile(client, file)
		} else {
			fmt.Printf("Creating file %s\n", file.Path)
			err = copyFile(client, file)
		}

		if err != nil {
			return err
		}
	}

	//Now, all files left on known should be deleted
	for path := range known {
		deleteFile(client, path)
	}

	return nil
}

// Deletes one file
func deleteFile(client *sshclient.Client, path string) error {
	fmt.Printf("Deleting file %s\n", path)
	cmd := fmt.Sprintf("sudo rm -f %s", path)
	out, err := client.Cmd(cmd).SmartOutput()
	if err != nil {
		log.Print(string(out))
		return err
	}
	return err
}

// Copies one file
// For simplicity, it will always override the remote files
func copyFile(client *sshclient.Client, file config.File) error {

	template := `
path="%s"
owner=%s
group=%s
mode=%d

cat <<EOF | sudo tee ${path}
%s
EOF

sudo chmod ${mode} ${path}
sudo chown ${owner}:${group} ${path}
`

	script := fmt.Sprintf(template,
		file.Path,
		file.Owner,
		file.Group,
		file.Mode,
		file.Content)
	out, err := client.Script(script).SmartOutput()
	if err != nil {
		log.Print(string(out))
		return err
	}

	return nil
}

// Starts the services based on which packages are services
// For simplicity, I'm not checking if service is running, just restarting them
// because it works regardless of the service is running or not
func startServices(client *sshclient.Client, app config.App) error {
	for _, p := range app.Packages {
		if p.IsService {
			script := fmt.Sprintf("sudo service %s restart", p.Name)
			out, err := client.Script(script).SmartOutput()
			if err != nil {
				log.Print(string(out))
				return err
			}
			fmt.Printf("Service %s started\n", p.Name)
		}
	}

	return nil
}

// Stops the services
func stopServices(client *sshclient.Client, app config.App) error {
	for _, p := range app.Packages {
		if p.IsService {
			script := fmt.Sprintf("sudo service %s stop", p.Name)
			out, err := client.Script(script).SmartOutput()
			if err != nil {
				log.Print(string(out))
				return err
			}
			fmt.Printf("Service %s stopped\n", p.Name)
		}
	}
	return nil
}

func aptUpdate(client *sshclient.Client) error {
	out, err := client.Cmd("sudo apt-get update -y").SmartOutput()
	if err != nil {
		log.Print(string(out))
		return err
	}
	return nil
}
