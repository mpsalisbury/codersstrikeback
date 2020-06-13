package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

// ****************************************************** Point

type Point struct {
	x, y float64
}

func (p *Point) plus(v *Vector) *Point {
	return &Point{p.x + v.vx, p.y + v.vy}
}

func (p1 *Point) minus(p2 *Point) *Vector {
	return &Vector{p1.x - p2.x, p1.y - p2.y}
}

func (p *Point) String() string {
	return fmt.Sprintf("%d %d", int(p.x), int(p.y))
}

func (p *Point) ReadInput(in io.Reader) {
	fmt.Fscan(in, &p.x, &p.y)
}

// ****************************************************** Vector

type Vector struct {
	vx, vy float64
}

func (v *Vector) norm() *Vector {
	len := v.len()
	return &Vector{v.vx / len, v.vy / len}
}

func (v1 *Vector) minus(v2 *Vector) *Vector {
	return &Vector{v1.vx - v2.vx, v1.vy - v2.vy}
}

func (v *Vector) perpendicular() *Vector {
	return &Vector{v.vy, -v.vx}
}

func (v1 *Vector) dot(v2 *Vector) float64 {
	return v1.vx*v2.vx + v1.vy*v2.vy
}

func (v *Vector) times(f float64) *Vector {
	return &Vector{v.vx * f, v.vy * f}
}

func (v *Vector) len() float64 {
	return math.Sqrt(v.len2())
}

func (v *Vector) len2() float64 {
	return v.dot(v)
}

func (v *Vector) String() string {
	return fmt.Sprintf("%d %d", int(v.vx), int(v.vy))
}

func (v *Vector) ReadInput(in io.Reader) {
	fmt.Fscan(in, &v.vx, &v.vy)
}

// ****************************************************** Action

type Action interface {
	String() string
}

type ThrustAction int

func (a ThrustAction) String() string {
	return fmt.Sprintf("%d", a)
}

type BoostAction struct{}

func (a BoostAction) String() string {
	return "BOOST"
}

type ShieldAction struct{}

func (a ShieldAction) String() string {
	return "SHIELD"
}

// ****************************************************** Pod

type Pod struct {
	pos              Point
	vel              Vector
	angle            float64
	nextCheckpointId int
	curLapNum        int

	strategy PodStrategy
	context  *Game
}

func (p *Pod) init(g *Game, s PodStrategy) {
	p.context = g
	p.strategy = s
	p.strategy.Init(p)
	p.curLapNum = -1
}

func (p1 *Pod) isAheadOf(p2 *Pod) bool {
	if p1.curLapNum > p2.curLapNum {
		return true
	}
	if p1.curLapNum < p2.curLapNum {
		return false
	}
	if p1.nextCheckpointId > p2.nextCheckpointId {
		return true
	}
	if p1.nextCheckpointId < p2.nextCheckpointId {
		return false
	}
	return p1.distToPoint(p1.nextCheckpoint(0)) <
		p2.distToPoint(p2.nextCheckpoint(0))
}

func (p *Pod) nextCheckpoint(offset int) *Point {
	index := p.nextCheckpointId + offset
	for index >= len(p.context.checkpoint) {
		index -= len(p.context.checkpoint)
	}
	return &p.context.checkpoint[index]
}

// Returns diff between ship angle and angle to given point in degrees.
func (p *Pod) angleToPoint(target *Point) float64 {
	dir := target.minus(&p.pos)
	pointAngle := math.Atan2(dir.vy, dir.vx) / math.Pi * 180.0
	diff := p.angle - pointAngle
	if diff < -180 {
		diff += 360
	} else if diff > 180 {
		diff -= 360
	}
	return diff
}

func (p *Pod) distToPoint(target *Point) float64 {
	return target.minus(&p.pos).len()
}

func (p *Pod) nextPos() *Point {
	return p.pos.plus(&p.vel)
}

// Returns the number of time steps at the current velocity until this pod is
// within the given checkpoint. Returns math.Inf(1) if it will not hit.
func (p *Pod) numStepsTo(cp *Point) float64 {
	r := 400.0
	cpToP := p.pos.minus(cp)
	a := p.vel.len2()
	b := 2 * p.vel.dot(cpToP)
	c := cpToP.len2() - r*r
	disc := b*b - 4*a*c
	if a == 0.0 || disc < 0.0 {
		return math.Inf(1)
	}
	return (-b + math.Sqrt(disc)) / (2 * a)
}

func (p1 *Pod) willCollideWith(p2 *Pod) bool {
	return p1.nextPos().minus(p2.nextPos()).len() < 800
}

func (p *Pod) ReadUpdateInput(in io.Reader) {
	p.pos.ReadInput(in)
	p.vel.ReadInput(in)
	oldCheckpointId := p.nextCheckpointId
	fmt.Fscan(in, &p.angle, &p.nextCheckpointId)
	if oldCheckpointId == 0 && p.nextCheckpointId == 1 {
		p.curLapNum++
	}
}

func (p *Pod) PrintActionOutput() {
	target, action := p.strategy.ComputeStep()
	fmt.Printf("%s %s\n", target.String(), action.String())
}

// ****************************************************** PodStrategy

type PodStrategy interface {
	Init(p *Pod)
	ComputeStep() (Point, Action)
}

// ****************************************************** LegacyStrategy
// Strategy used in earlier stages.
// Go to next checkpoint directly.
// Boost when heading straight towards a checkpoint.
// Shield if opponent pod is about to collide.

type LegacyStrategy struct {
	p *Pod

	dontBoostBeforeLapNum int
	usedBoost             bool

	target Point
	action Action
}

func (s *LegacyStrategy) Init(p *Pod) {
	s.p = p
}

func (s *LegacyStrategy) ComputeStep() (Point, Action) {
	s.target = s.ComputeTarget()
	s.action = s.ComputeAction()
	return s.target, s.action
}

func (s *LegacyStrategy) ComputeTarget() Point {
	p := s.p
	nextCheckpointVector := p.nextCheckpoint(0).minus(&p.pos)
	perpCheckpointVector := nextCheckpointVector.norm().perpendicular()
	offsetLen := p.vel.dot(perpCheckpointVector)

	return *p.nextCheckpoint(0).plus(perpCheckpointVector.times(-2 * offsetLen))
}

func (s *LegacyStrategy) ComputeAction() Action {
	p := s.p
	if p.willCollideWith(&p.context.oppPod[0]) ||
		p.willCollideWith(&p.context.oppPod[1]) {
		return ShieldAction{}
	}

	targetDist := s.target.minus(&p.pos).len()
	angleDiff := math.Abs(p.angleToPoint(&s.target))
	if !s.usedBoost && s.dontBoostBeforeLapNum >= p.curLapNum &&
		angleDiff < 10 && targetDist > 4000 {
		s.usedBoost = true
		return BoostAction{}
	}

	thrust := ThrustAction(0)
	if angleDiff < 90.0 {
		thrust = 100
	}
	if p.distToPoint(p.nextCheckpoint(0)) < 2000 {
		thrust /= 2
	}
	return thrust
}

// ****************************************************** FlyStrategy
// Go as fast as possible through the checkpoints with
// no regard for other pods.
type FlyStrategy struct {
	p         *Pod
	usedBoost bool
}

func (s *FlyStrategy) Init(p *Pod) {
	s.p = p
}

func (s *FlyStrategy) ComputeStep() (Point, Action) {
	p := s.p
	if p.numStepsTo(p.nextCheckpoint(0)) < 6 {
		return s.DriftStep()
	} else {
		return s.BlastStep()
	}
}

func (s *FlyStrategy) TargetForGoal(goal *Point) *Point {
	p := s.p
	targetVector := goal.minus(&p.pos)
	perpVector := targetVector.norm().perpendicular()
	offsetLen := p.vel.dot(perpVector)
	return goal.plus(perpVector.times(-3 * offsetLen))
}

func (s *FlyStrategy) ShouldShieldAgainst(other *Pod) bool {
	return s.p.willCollideWith(other) &&
		s.p.vel.minus(&other.vel).len() > 100
}

func (s *FlyStrategy) ShouldShield() bool {
	return s.ShouldShieldAgainst(&s.p.context.oppPod[0]) ||
		s.ShouldShieldAgainst(&s.p.context.oppPod[1])
}

// Go as fast as possible to the next checkpoint.
func (s *FlyStrategy) BlastStep() (Point, Action) {
	p := s.p
	target := s.TargetForGoal(p.nextCheckpoint(0))

	if s.ShouldShield() {
		return *target, ShieldAction{}
	}

	targetDist := target.minus(&p.pos).len()
	angleDiff := math.Abs(p.angleToPoint(target))
	if !s.usedBoost && angleDiff < 10 && targetDist > 4000 {
		s.usedBoost = true
		return *target, BoostAction{}
	}

	thrust := ThrustAction(0)
	if angleDiff < 90.0 {
		thrust = 100
	}
	return *target, thrust
}

// Assume that we'll drift into the checkpoint and start turning to the next one.
func (s *FlyStrategy) DriftStep() (Point, Action) {
	target := s.TargetForGoal(s.p.nextCheckpoint(1))
	if s.ShouldShield() {
		return *target, ShieldAction{}
	}

	angleDiff := math.Abs(s.p.angleToPoint(target))
	thrust := ThrustAction(0)
	if angleDiff < 45.0 {
		thrust = 100
	}
	return *target, thrust
}

// ****************************************************** BlockStrategy
// Block the lead opponent pod from getting to a future checkpoint.
type BlockStrategy struct {
	p          *Pod
	blockedPod *Pod
}

func (s *BlockStrategy) Init(p *Pod) {
	s.p = p
	s.blockedPod = &p.context.oppPod[0]
}

func (s *BlockStrategy) ComputeStep() (Point, Action) {
	/*
	   if blocked pod has passed, choose new blocked pod
	   go to blocking position
	   aim at pod
	   when close enough, blast towards it and shield
	   try to stay between pod and its goal
	*/
	target := s.ComputeTarget()
	return *target, ThrustAction(100)
}

func (s *BlockStrategy) ComputeTarget() *Point {
	speed := math.Max(100, s.p.vel.len())
	t, ok := s.ComputeHitTime(speed)
	if !ok {
		return &s.blockedPod.pos
	}
	return s.blockedPod.pos.plus(s.blockedPod.vel.times(t))
}

func (s *BlockStrategy) ComputeHitTime(speed float64) (float64, bool) {
	v2 := s.blockedPod.vel
	delta := s.blockedPod.pos.minus(&s.p.pos)
	a := v2.len2() - speed*speed
	b := 2 * delta.dot(&v2)
	c := delta.len2()
	disc := b*b - 4*a*c
	if a == 0 || disc < 0 {
		return 0, false
	} else {
		t := (-b - math.Sqrt(disc)) / (2 * a)
		if t < 0 {
			t = (-b + math.Sqrt(disc)) / (2 * a)
		}
		if t < 0 {
			return 0, false
		}
		return t, true
	}
}

// ****************************************************** Game

type Game struct {
	myPod      [2]Pod
	oppPod     [2]Pod
	checkpoint []Point

	numLaps int
}

func (g *Game) Init(in io.Reader) {
	g.ReadSetupInput(in)
	g.myPod[0].init(g, &FlyStrategy{})
	g.myPod[1].init(g, &BlockStrategy{})
}

func (g *Game) ReadSetupInput(in io.Reader) {
	fmt.Scan(&g.numLaps)
	var numCheckpoints int
	fmt.Scan(&numCheckpoints)
	for i := 0; i < numCheckpoints; i++ {
		var checkpoint Point
		checkpoint.ReadInput(in)
		g.checkpoint = append(g.checkpoint, checkpoint)
	}
}

func (g *Game) ReadPodUpdateInput(in io.Reader) {
	g.myPod[0].ReadUpdateInput(in)
	g.myPod[1].ReadUpdateInput(in)
	g.oppPod[0].ReadUpdateInput(in)
	g.oppPod[1].ReadUpdateInput(in)
}

func (g *Game) PrintPodActionOutput() {
	g.myPod[0].PrintActionOutput()
	g.myPod[1].PrintActionOutput()
}

// ****************************************************** main

func main() {
	var g Game
	in := os.Stdin
	g.Init(in)
	for {
		g.ReadPodUpdateInput(in)
		g.PrintPodActionOutput()
	}
}
