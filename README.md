# Coldline Miami

A top-down action game inspired by Hotline Miami, built in Go using [Raylib](https://github.com/gen2brain/raylib-go).

## Requirements

- Go 1.26+
- Raylib system dependencies (see below)

### macOS

```sh
xcode-select --install
brew install raylib
```

### Linux

```sh
sudo apt install libraylib-dev
# or
sudo pacman -S raylib
```

## Build & Run

```sh
./build.sh
```

This cleans the `build/` directory, compiles, and launches the game in one step. The binary is output to `build/coldline-miami`.

To build without running:

```sh
go build -o build/coldline-miami .
```

## Controls

| Key / Input | Action |
|-------------|--------|
| Arrow keys | Move player |
| Space | Shoot toward mouse cursor |

## Architecture

The game uses a simple **Entity Component System (ECS)** pattern. All game state lives in a `World`, which holds separate maps from entity ID to each component type. Systems operate on those maps each frame.

```
World
├── entities        — all active entities (UUID → Entity)
├── transforms      — position (X, Y)
├── movements       — player-controlled movement speed
├── staticMovements — direction + speed for autonomous movers (bullets)
├── healths         — current/max HP and invincibility timer
├── colliders       — AABB hitbox + tag (player/enemy/wall/bullet)
├── shooters        — fire rate cooldown
└── dimensions      — logical resolution (width × height) used for scaling and despawn bounds
```

### Systems (update order)

| System | What it does |
|--------|-------------|
| `DespawnSystem` | Kills entities that have travelled outside the logical world bounds |
| `MovementSystem` | Reads arrow key input, moves entities with a `Movement` component |
| `ShootingSystem` | Spawns a bullet aimed at the mouse cursor on Space, respects fire rate cooldown |
| `StaticMovementSystem` | Moves entities along a fixed normalised direction vector (bullets) |
| `CollisionSystem` | AABB overlap checks across all collidables, dispatches to `onCollision` |
| `HealthSystem` | Ticks down invincibility timers, kills entities at 0 HP |
| `GameOverSystem` | Sets `endState` to Victory or Defeat when player/enemies are gone |

### Collision rules (`onCollision`)

| Pair | Result |
|------|--------|
| Bullet + Enemy | Enemy takes 1 damage, bullet destroyed |
| Player + Enemy | Player takes 1 damage (with 1s invincibility window) |
| Player/Enemy + Wall | No-op (wall push-back not yet implemented) |

### Adding a new entity

```go
e := world.NewEntity("MyThing")
world.AddTransform(e, 100, 200)
world.AddHealth(e, 3, 3)
world.AddCollider(e, 16, 16, core.TagEnemy)
```

### Adding a new component

1. Define the struct in [src/core/components.go](src/core/components.go)
2. Add a `map[uuid.UUID]*YourComponent` field to `World` in [src/core/world.go](src/core/world.go)
3. Initialise the map in `NewWorld()`
4. Add an `AddYourComponent` helper on `World`
5. Add `delete(w.yourComponents, id)` to `KillEntity`
6. Write a system function in [src/core/systems.go](src/core/systems.go) and call it from `update()` in [src/core/loop.go](src/core/loop.go)

## Debug Mode

Debug mode is enabled by default (`world.TurnOnDebug()` in `init`). It overlays:

- "Debug Mode ON" label
- Current FPS
- HP value above each entity that has a `Health` component

Toggle with `world.TurnOffDebug()`.

## Rendering

The game renders at a fixed **logical resolution** (800×600) into a `RenderTexture2D`, then scales that texture to fill whatever the actual window size is. This means all gameplay coordinates, collision, and spawning logic works in logical pixels — the window can be resized freely without touching any game code.

Mouse input is remapped to logical coordinates via `getLogicalMousePosition` in [src/core/utils.go](src/core/utils.go) so aiming stays accurate at any window size.

## Project Structure

```
coldline-miami/
├── main.go              — entry point
├── build.sh             — build + run script
├── go.mod / go.sum      — module definition and locked checksums
├── build/               — compiled binaries (not committed)
└── src/
    └── core/
        ├── components.go — component structs and tag constants
        ├── world.go      — World struct, entity factory, component helpers
        ├── systems.go    — all ECS systems
        ├── loop.go       — game loop, init, render pipeline
        ├── game.go       — Game struct (holds render texture)
        ├── entity.go     — Entity struct
        ├── settings.go   — debug/settings flags
        └── utils.go      — shared helpers (logical mouse position, etc.)
```
