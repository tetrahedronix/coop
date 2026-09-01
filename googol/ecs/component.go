package ecs

type ComponentTypeID uint64

const (
	ComponentTypePosition ComponentTypeID = 1 << iota
	ComponentTypeSelectable
	ComponentTypeShape
	ComponentTypeVelocity
	ComponentTypeTile
)

// Components have data
type Component interface {
	TypeID() ComponentTypeID
	Add(interface{})
	Copy(interface{})
	Get() interface{}
	Reset()
}

// const (
// 	haveRender ComponentTypeID = 1 << 0
// 	havePhysic ComponentTypeID = 1 << 1
// 	haveInput  ComponentTypeID = 1 << 2
// 	haveHealth ComponentTypeID = 1 << 3
// 	haveSound  ComponentTypeID = 1 << 4
// )

func CopyComponent(src Component) Component {
	switch src := src.(type) {
	case *Position:
		dst := NewPosition().(*Position)
		dst.coordinate = src.coordinate
		return dst
	case *Selectable:
		dst := NewSelectable().(*Selectable)
		dst.selected = src.selected
		return dst
	case *Shape:
		dst := NewShape().(*Shape)
		dst.primitive = src.primitive
		return dst
	case *Velocity:
		dst := NewVelocity().(*Velocity)
		dst.speed = src.speed
		dst.direction = src.direction
		return dst

		// Altri componenti
		// ...

	default:
		panic("unsupported component type")
	}
}
