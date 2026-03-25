package core

const DESPAWN_OFFSET = 100.0
const PLAYER_INVICIBLE_AFTER_HIT_TIME = 1.0

const WORLD_WIDTH = 800.0
const WORLD_HEIGHT = 600.0

const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = 720

const DEBUG_MODEON = false

type GameState string

const (
	StatePlaying  GameState = "playing"
	StatePaused   GameState = "paused"
	StateGameOver GameState = "game_over"
)

type GameEndState string

const (
	Victory     GameEndState = "victory"
	Defeat      GameEndState = "defeat"
	In_Progress GameEndState = "in_progress"
)

type ColliderTag string

const (
	TagPlayer ColliderTag = "player"
	TagEnemy  ColliderTag = "enemy"
	TagWall   ColliderTag = "wall"
	TagBullet ColliderTag = "bullet"
)
