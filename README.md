# size-server

used to investigate https://github.com/traefik/traefik/issues/10687

The server only has one handler at `/`. It accepts a parameter `size` that
determines the size of the response payload.

By default the program runs as a server. There is a `-client` flag that can
switch the program to run as a client of the server instead.

## Testing
One instance running as a server behind traefik with the config in the 
docker-compose.yml and a separate instance running on the host as a client
making increasingly larger requests to understand the behavior of the reverse
proxy when the response size crosses the boundary set by the `maxResponseBodyBytes`
of the proxy.

## Logs
The logs named `processed-log-{number}.txt` are named with the value of the
`maxResponseBodyBytes` that was used when collecting the given log, e.g.
`processed-log-2000.txt` means `maxResponseBodyBytes` was set to 2000 at the time.