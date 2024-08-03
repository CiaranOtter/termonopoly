package game

type OwnerInterface interface {
	OfferProperty(SpaceInterface)
	ChargeRent(int)
	PassGo()
	GoToJail()
}
