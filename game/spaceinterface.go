package game

type SpaceInterface interface {
	OnLand(OwnerInterface)
	OnPass(OwnerInterface)
	SetNext(SpaceInterface)
	SetPrev(SpaceInterface)
	Print()
	GetNext() SpaceInterface
	GetPrev() SpaceInterface

	Afford(int) (int, bool)
	GetGroup() string
	SetGroup(string)
	SetOwner(OwnerInterface)
}
