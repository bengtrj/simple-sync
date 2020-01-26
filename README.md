**Simple configuration management tool:**
- Problem states that the tool should be rudimentary
- Problem suggests spending ~4 hours including other exercise
- asks for a way to install/remove packages
- asks for a way of defining files and file metadata
- asks for idempotency

To me this meant keeping it simple, focus on getting a `php` page up, in a repeatable way, on a list of remote hosts I happen to have access. If you run the tool repeateadly with the same configuration, it doesn't change the effective state of the system and the server is still running.

For the scope of this exercise, I think a good start would be:
- Start by writing a tool that can:
    - ssh into a host
    - install apache2
    - try to access it

