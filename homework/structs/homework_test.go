package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) Option {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) Option {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) Option {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) Option {
	return func(person *GamePerson) {
		person.mana = uint16(mana)
	}
}

func WithHealth(health int) Option {
	return func(person *GamePerson) {
		person.health = uint16(health)
	}
}

func WithRespect(respect int) Option {
	return func(person *GamePerson) {
		person.setAttribute(0, 4, uint32(respect))
	}
}

func WithStrength(strength int) Option {
	return func(person *GamePerson) {
		person.setAttribute(4, 4, uint32(strength))
	}
}

func WithExperience(experience int) Option {
	return func(person *GamePerson) {
		person.setAttribute(8, 4, uint32(experience))
	}
}

func WithLevel(level int) Option {
	return func(person *GamePerson) {
		person.setAttribute(12, 4, uint32(level))
	}
}

func WithHouse() Option {
	return func(person *GamePerson) {
		person.setAttribute(16, 1, 1)
	}
}

func WithGun() Option {
	return func(person *GamePerson) {
		person.setAttribute(17, 1, 1)
	}
}

func WithFamily() Option {
	return func(person *GamePerson) {
		person.setAttribute(18, 1, 1)
	}
}

func WithType(personType int) Option {
	return func(person *GamePerson) {
		person.setAttribute(19, 2, uint32(personType))
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name       [40]byte
	x          int32
	y          int32
	z          int32
	gold       int32
	mana       uint16
	health     uint16
	attributes uint32
}

func (p *GamePerson) setAttribute(offset, size int, value uint32) {
	mask := ((uint32(1) << size) - 1) << offset
	p.attributes &= ^mask
	p.attributes |= (value << offset) & mask
}

func (p *GamePerson) getAttribute(offset, size int) uint32 {
	mask := uint32(((1 << size) - 1) << offset)
	return (p.attributes & mask) >> offset
}

func NewGamePerson(options ...Option) GamePerson {
	var person GamePerson
	for _, opt := range options {
		opt(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	name := p.name[:]
	for i, b := range name {
		if b == 0 {
			name = name[:i]
			break
		}
	}
	return string(name)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.mana)
}

func (p *GamePerson) Health() int {
	return int(p.health)
}

func (p *GamePerson) Respect() int {
	return int(p.getAttribute(0, 4))
}

func (p *GamePerson) Strength() int {
	return int(p.getAttribute(4, 4))
}

func (p *GamePerson) Experience() int {
	return int(p.getAttribute(8, 4))
}

func (p *GamePerson) Level() int {
	return int(p.getAttribute(12, 4))
}

func (p *GamePerson) HasHouse() bool {
	return p.getAttribute(16, 1) == 1
}

func (p *GamePerson) HasGun() bool {
	return p.getAttribute(17, 1) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return p.getAttribute(18, 1) == 1
}

func (p *GamePerson) Type() int {
	return int(p.getAttribute(19, 2))
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
