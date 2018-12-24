package cassgowary

import (
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"
)

const (
	LEFT    = "left"
	RIGHT   = "right"
	TOP     = "top"
	BOTTOM  = "bottom"
	HEIGHT  = "height"
	WIDTH   = "width"
	CENTERX = "centerX"
	CENTERY = "centerY"
)

var constraints = []string{
	"container.columnWidth == container.width * 0.4",
	"container.thumbHeight == container.columnWidth / 2",
	"container.padding == container.width * (0.2 / 3)",
	"container.leftPadding == container.padding",
	"container.rightPadding == container.width - container.padding",
	"container.paddingUnderThumb == 5",
	"container.rowPadding == 15",
	"container.buttonPadding == 20",

	"thumb0.left == container.leftPadding",
	"thumb0.top == container.padding",
	"thumb0.height == container.thumbHeight",
	"thumb0.width == container.columnWidth",

	"title0.left == container.leftPadding",
	"title0.top == thumb0.bottom + container.paddingUnderThumb",
	"title0.height == title0.intrinsicHeight",
	"title0.width == container.columnWidth",

	"thumb1.right == container.rightPadding",
	"thumb1.top == container.padding",
	"thumb1.height == container.thumbHeight",
	"thumb1.width == container.columnWidth",

	"title1.right == container.rightPadding",
	"title1.top == thumb0.bottom + container.paddingUnderThumb",
	"title1.height == title1.intrinsicHeight",
	"title1.width == container.columnWidth",

	"thumb2.left == container.leftPadding",
	"thumb2.top >= title0.bottom + container.rowPadding",
	"thumb2.top == title0.bottom + container.rowPadding !weak",
	"thumb2.top >= title1.bottom + container.rowPadding",
	"thumb2.top == title1.bottom + container.rowPadding !weak",
	"thumb2.height == container.thumbHeight",
	"thumb2.width == container.columnWidth",

	"title2.left == container.leftPadding",
	"title2.top == thumb2.bottom + container.paddingUnderThumb",
	"title2.height == title2.intrinsicHeight",
	"title2.width == container.columnWidth",

	"thumb3.right == container.rightPadding",
	"thumb3.top == thumb2.top",

	"thumb3.height == container.thumbHeight",
	"thumb3.width == container.columnWidth",

	"title3.right == container.rightPadding",
	"title3.top == thumb3.bottom + container.paddingUnderThumb",
	"title3.height == title3.intrinsicHeight",
	"title3.width == container.columnWidth",

	"thumb4.left == container.leftPadding",
	"thumb4.top >= title2.bottom + container.rowPadding",
	"thumb4.top >= title3.bottom + container.rowPadding",
	"thumb4.top == title2.bottom + container.rowPadding !weak",
	"thumb4.top == title3.bottom + container.rowPadding !weak",
	"thumb4.height == container.thumbHeight",
	"thumb4.width == container.columnWidth",

	"title4.left == container.leftPadding",
	"title4.top == thumb4.bottom + container.paddingUnderThumb",
	"title4.height == title4.intrinsicHeight",
	"title4.width == container.columnWidth",

	"thumb5.right == container.rightPadding",
	"thumb5.top == thumb4.top",
	"thumb5.height == container.thumbHeight",
	"thumb5.width == container.columnWidth",

	"title5.right == container.rightPadding",
	"title5.top == thumb5.bottom + container.paddingUnderThumb",
	"title5.height == title5.intrinsicHeight",
	"title5.width == container.columnWidth",

	"line.height == 1",
	"line.width == container.width",
	"line.top >= title4.bottom + container.rowPadding",
	"line.top >= title5.bottom + container.rowPadding",

	"more.top == line.bottom + container.buttonPadding",
	"more.height == more.intrinsicHeight",
	"more.left == container.leftPadding",
	"more.right == container.rightPadding",

	"container.height == more.bottom + container.buttonPadding",
}

type nodeMap map[string]*Variable
type nodesMap map[string]nodeMap
type gridVariableResolver struct {
	solver *Solver
	nodes  nodesMap
}

func (vr *gridVariableResolver) getVariableFromNode(node nodeMap, variableName string) *Variable {
	if v, exists := node[variableName]; exists {
		return v
	}
	v := NewVariable(variableName)
	node[variableName] = v

	switch variableName {
	case RIGHT:
		v2 := vr.getVariableFromNode(node, LEFT)
		v3 := vr.getVariableFromNode(node, WIDTH)
		e := v2.Add(v3)
		c := v.EqualsExpression(e)
		vr.solver.AddConstraint(c)
	case BOTTOM:
		v2 := vr.getVariableFromNode(node, TOP)
		v3 := vr.getVariableFromNode(node, HEIGHT)
		e := v2.Add(v3)
		c := v.EqualsExpression(e)
		vr.solver.AddConstraint(c)
	}
	return v
}

func (vr *gridVariableResolver) getNode(nodeName string) map[string]*Variable {
	if node, exists := vr.nodes[nodeName]; exists {
		return node
	}

	node := map[string]*Variable{
		nodeName: nil,
	}
	vr.nodes[nodeName] = node
	return node
}

func (vr *gridVariableResolver) ResolveVariable(name string) (*Variable, error) {
	arr := strings.Split(name, ".")
	if len(arr) == 2 {
		nodeName, propertyName := arr[0], arr[1]
		node := vr.getNode(nodeName)
		v := vr.getVariableFromNode(node, propertyName)
		return v, nil
	}
	return nil, errors.New("can't resolve variable")
}

func (vr *gridVariableResolver) ResolveConstant(name string) (*Expression, error) {

	f, err := strconv.ParseFloat(name, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse '%s'", name)
	}
	e := NewExpression(f)
	return e, nil
}

func createGridVariableResolver(solver *Solver, nodes nodesMap) VariableResolver {
	return &gridVariableResolver{solver, nodes}
}

func TestGridLayout(t *testing.T) {
	solver := NewSolver()
	nodes := nodesMap{}
	variableResolver := createGridVariableResolver(solver, nodes)

	cp := NewConstraintParser()
	for _, constraint := range constraints {
		c, err := cp.ParseConstraint(constraint, variableResolver)
		assert.NoError(t, err)
		if err != nil {
			return
		}
		err = solver.AddConstraint(c)
		assert.NoError(t, err)
		if err != nil {
			return
		}
	}

	c, err := cp.ParseConstraint("container.width == 300", variableResolver)
	assert.NoError(t, err)
	solver.AddConstraint(c)
	c, err = cp.ParseConstraint("title0.intrinsicHeight == 100", variableResolver)
	assert.NoError(t, err)
	solver.AddConstraint(c)
	c, err = cp.ParseConstraint("title1.intrinsicHeight == 110", variableResolver)
	assert.NoError(t, err)
	solver.AddConstraint(c)
	c, err = cp.ParseConstraint("title2.intrinsicHeight == 120", variableResolver)
	assert.NoError(t, err)
	solver.AddConstraint(c)
	c, err = cp.ParseConstraint("title3.intrinsicHeight == 130", variableResolver)
	assert.NoError(t, err)
	solver.AddConstraint(c)
	c, err = cp.ParseConstraint("title4.intrinsicHeight == 140", variableResolver)
	assert.NoError(t, err)
	solver.AddConstraint(c)
	c, err = cp.ParseConstraint("title5.intrinsicHeight == 150", variableResolver)
	assert.NoError(t, err)
	solver.AddConstraint(c)
	c, err = cp.ParseConstraint("more.intrinsicHeight == 160", variableResolver)
	assert.NoError(t, err)
	solver.AddConstraint(c)

	solver.UpdateVariables()

	assert.InDelta(t, 20, nodes["thumb0"]["top"].Value, Epsilon)
	assert.InDelta(t, 20, nodes["thumb1"]["top"].Value, Epsilon)
	assert.InDelta(t, 85, nodes["title0"]["top"].Value, Epsilon)
	assert.InDelta(t, 85, nodes["title1"]["top"].Value, Epsilon)
	assert.InDelta(t, 200, nodes["thumb2"]["top"].Value, Epsilon)
	assert.InDelta(t, 200, nodes["thumb3"]["top"].Value, Epsilon)
	assert.InDelta(t, 265, nodes["title2"]["top"].Value, Epsilon)
	assert.InDelta(t, 265, nodes["title3"]["top"].Value, Epsilon)
	assert.InDelta(t, 400, nodes["thumb4"]["top"].Value, Epsilon)
	assert.InDelta(t, 400, nodes["thumb5"]["top"].Value, Epsilon)
	assert.InDelta(t, 465, nodes["title4"]["top"].Value, Epsilon)
	assert.InDelta(t, 465, nodes["title5"]["top"].Value, Epsilon)
}

// func TestGridX1000(t *testing.T) {
// 	start := time.Now()
// 	for i := 0; i < 1000; i++ {
// 		TestGridLayout(t)
// 	}
// 	log.Printf("testGridX1000 took %s.", time.Since(start))
// }

func printNodes(variables nodesMap) {
	for nodeName, nodes := range variables {
		log.Printf("node: %s", nodeName)
		printVariables(nodes)
	}
}

func printVariables(nodes nodeMap) {
	for name, v := range nodes {

		log.Printf(" %s = %f (address:%+v)", name, v.Value, v)
	}
}
