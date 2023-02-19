package room

type Vote interface {
	Value() float32
	Type() string
}

type VoteTemplate struct {
	Title string
	Votes []Vote
}

func NewVoteTemplate(title string, votes []Vote) VoteTemplate {
	return VoteTemplate{Title: title, Votes: votes}
}

type VoteNumeric struct {
	value float32
}

func NewVoteNumeric(value float32) VoteNumeric {
	return VoteNumeric{value: value}
}

func (v VoteNumeric) Value() float32 {
	return v.value
}

func (v VoteNumeric) Type() string {
	return "value"
}

type VoteSpecial struct {
	t string
}

func NewVoteSpecial(t string) VoteSpecial {
	return VoteSpecial{t: t}
}

func (v VoteSpecial) Value() float32 {
	return 0
}

func (v VoteSpecial) Type() string {
	return v.t
}

var VoteUnknown = NewVoteSpecial("unknown")
var VoteBreak = NewVoteSpecial("break")
var VoteInfinity = NewVoteSpecial("infinity")

var VoteTemplateFibonacci = NewVoteTemplate("Fibonacci", []Vote{
	VoteUnknown,
	VoteBreak,
	NewVoteNumeric(0),
	NewVoteNumeric(1),
	NewVoteNumeric(2),
	NewVoteNumeric(3),
	NewVoteNumeric(5),
	NewVoteNumeric(8),
	NewVoteNumeric(13),
	NewVoteNumeric(21),
	NewVoteNumeric(34),
	NewVoteNumeric(55),
	NewVoteNumeric(89),
	VoteInfinity,
})

var VoteTemplateModFibonacci = NewVoteTemplate("Modified fibonacci", []Vote{
	VoteUnknown,
	VoteBreak,
	NewVoteNumeric(0),
	NewVoteNumeric(0.5),
	NewVoteNumeric(1),
	NewVoteNumeric(2),
	NewVoteNumeric(3),
	NewVoteNumeric(5),
	NewVoteNumeric(8),
	NewVoteNumeric(13),
	NewVoteNumeric(20),
	NewVoteNumeric(40),
	NewVoteNumeric(100),
	VoteInfinity,
})

var DefaultVoteTemplates = []VoteTemplate{
	VoteTemplateFibonacci,
	VoteTemplateModFibonacci,
}
