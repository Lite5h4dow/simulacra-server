# Simulacra
A Protocol to share your world with the rest of the world

## Mission Statement
Anyones idea, hosted publically, no reliance on a single client (VRC, Resonite, Neos)
Re-democratise the virtual space for everyone using the same tools that made the world wide web accessible
Simulacra aims to make creating a world as simple as making it in blender and (eventually) just clicking build.
assets, materials, and transforms all baked into a collection of assets to be streamed just like a website.

Modern day tech is already capable of rendering scenes even in our web browsers, so thats where i'll start.
a webserver with a hosted webclient that lets you jump straight into the world, no installing anything.
eventually i want to make a dedicated client using godot that will have extra benefits i wont be able to
replicate with a browser, like LOD and more complex render tech, but im starting small with the aim to just
render a cube in vr from a simulacra server from a manifest.

## Todo:
- [x] Consume TCP/IP Packet
- [x] Detect HTTP / Simulacra
- [x] Handle HTTP request
- [x] Handle Simulacra request
- [x] Respond with Manifest
- [x] Respond with Asset
- [x] Respond with Placeholder HTTP response
- [ ] Make Webclient (Starting with A-Frame for a simple POC)
- [ ] Add Materials to manifest
- [ ] Respond with Material (possibly another manifest or an archive with all material assets)
- [ ] Compression (?) (Need to test decompression efficiency in browser)

## Nice-To-Haves:
- [-] Flags: 
  - [x] Custom Port
  - [x] Custom Directory
  - [ ] Help Manual
  - [ ] (I'm sure ill think of more flags)

- [ ] nix-shell for development (make it easier for others to contribute).
- [ ] nix-flake for CI/CD (Binary and OCI Container).
- [ ] CI/CD Pipeline
- [ ] Auto Pull / Update client from base server configs / flags.
- [ ] Systemd Service setup.
- [ ] P2P multiplayer
  - [ ] will TCP be fast enough?
  - [ ] will UDP Holepunching be required?
  - [ ] Support Both? (Probably not)
  - [ ] P2P websockets maybe?

## Longshots:
- [ ] Custom Avatars (will require P2P first).
- [ ] Face Tracking.
- [ ] Eye Tracking.
- [ ] Full Body Tracking.
- [ ] Blender Addon (including custom extra primitives to represent spawn points, portals, and other Simulacra specific entities).
- [ ] Client Libraries:
  - [ ] Go (custom clients and bots/scripts maybe?).
  - [ ] C# (Godo-Mono, Resonite Cross compat?, Unity).
    - This will be used for the official dedicated client (if Godot-mono can cross compile for quest/pico/vive).
  - [ ] JS/TS (web clients / websites).
  - [ ] GDScript (?).
    - I will not use this (probably, unless godot-mono is still missing features by the time i get to it).
    - will only add this to the scope if there is demand for it.
  - [ ] C++ ( Unreal engine ??).
    - is this even worth it? UE4/5 is so heavy for standalone headsets.
    - will only build if:
      1. I (or a community member) feel(s) confident writing it efficiently.
      2. I (or a community member) has time.

## Links:
  - Site: https://litelot.us (currently down, home server is down for maintenance)
  - Blog: https://litelot.us/blog (currently down, home server is down for maintenance)
  - Discord: https://discord.gg/pasEV9p4 (please dont blow this up too much)
