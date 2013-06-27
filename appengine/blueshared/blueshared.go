package blueshared


type Get10Options struct {
  NoAttack               bool
  NoAttackWithoutDefense bool
  SetsAvailable          []int
  ForbiddenCards         []int
}
