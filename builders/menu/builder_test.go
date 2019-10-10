package menu

import (
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type menuBuilderTestSuite struct {
	suite.Suite
}

func TestMenuBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(menuBuilderTestSuite))
}

func (s *menuBuilderTestSuite) TestMenuBuilderTestSuite_CreateNew() {

	fn := func(ctx types.MessageContextInterface) {}

	mb := NewMenuBuilder(
		SelectiveKeyboard(true),
		ResizeKeyboard(true),
		OneTimeKeyboard(false))

	mb.Button("menu1", fn).
		Button("menu2", fn).NewRow().
		Button("row2", fn)

	mm := mb.Build()

	s.Require().True(mm.Contain("menu1"), "menu1 must is contain")

}
