package main

type Ascii [3][3]rune
type Ascii_Display = [2]int

const Container_Size = 7
const Ascii_Width = 3
const Ascii_Height = 3
const Border_Width = Container_Size + 4
const Border_Height = Container_Size + 1
const Gap = 1

var zero = Ascii{
	{'█', '▀', '█'},
	{'█', ' ', '█'},
	{'▀', '▀', '▀'},
}
var one = Ascii{
	{'▝', '█', ' '},
	{' ', '█', ' '},
	{'▀', '▀', '▀'},
}
var two = Ascii{
	{'▀', '▀', '█'},
	{'█', '▀', '▀'},
	{'▀', '▀', '▀'},
}
var three = Ascii{
	{'▀', '▀', '█'},
	{'▀', '▀', '█'},
	{'▀', '▀', '▀'},
}
var four = Ascii{
	{'█', ' ', '█'},
	{'▀', '▀', '█'},
	{' ', ' ', '▀'},
}
var five = Ascii{
	{'█', '▀', '▀'},
	{'▀', '▀', '█'},
	{'▀', '▀', '▀'},
}
var six = Ascii{
	{'█', '▀', '▀'},
	{'█', '▀', '█'},
	{'▀', '▀', '▀'},
}
var seven = Ascii{
	{'▀', '▀', '█'},
	{' ', ' ', '█'},
	{' ', ' ', '▀'},
}
var eight = Ascii{
	{'█', '▀', '█'},
	{'█', '▀', '█'},
	{'▀', '▀', '▀'},
}
var nine = Ascii{
	{'█', '▀', '█'},
	{'▀', '▀', '█'},
	{'▀', '▀', '▀'},
}
var A = Ascii{
	{'█', '▀', '█'},
	{'█', '▀', '█'},
	{'▀', ' ', '▀'},
}
var I = Ascii{
	{'▀', '█', '▀'},
	{' ', '█', ' '},
	{'▀', '▀', '▀'},
}
var P = Ascii{
	{'█', '▀', '█'},
	{'█', '▀', '▀'},
	{'▀', ' ', ' '},
}
