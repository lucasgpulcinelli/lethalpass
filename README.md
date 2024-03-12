# Lethal Pass
This program is a simple adapter from an IPV6 UDP external connection to a IPV4 UDP listener.

The code was written specifically to allow a unity game host to use IPV6 as the public address for a game session.

## The problem
Game programmers are not network experts, nor are modders. In some games that have non steam play like modes, the host and guests cannot set the port that the game is hosted on. This, as well as the fact that games often use UDP, and that some only accept ip addresses, make it so free solutions such as ngrok or localtunnel cannot be used to make host session avaliable over the internet to the guests (because either they give you a random port or a temporary hostname, or only accept tcp).

Also, users are kinda dumb. I didn't want anyone to have to install something like a shared tailscale session in order to create a virtual LAN. 

## The solution (or the best I could find at least): IPV6!
IPV6 provides a individual public IP address that can be interpreted by unity for clients. So if both the host and guests have an IPV6 enabled network we are ok right? RIGHT?

There is a small problem: It seems like at least the game I was interested in supports IPV6 guests but only exposes the game for the host in IPV4... this is where lethal pass comes in! We can define a UDP6 listener, and, when an external guest connects, create an UPD4 connection and forward the external packets to the real game UDP4 listener. This needs some minor state keeping to associate the external connection to the local one, but this is mostly painless.

## How to run lethalpass
Just use `go run .` in the root of the project, there are no dependencies aside from the go standard library.
