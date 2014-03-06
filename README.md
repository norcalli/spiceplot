spiceplot
=========

A program written in Go and using plotinum that interprets raw spice files and produces pretty plots in __.svg, .png, .jpg, .jpeg, .eps, .tiff, .pdf__ formats.

##Features:
- Plotting specific variables
- Calculating expressions

### Future:
- Plotting expressions
- Complex variables
- Calculating expressions better

#Usage:

```
Usage: ./spiceplot [flags] <rawfile> <output>

Supported formats for the output are: .svg, .png, .jpg, .jpeg, .eps, .tiff, .pdf.
Possible flags:
  -e="": Comma separated string of expressions to evaluate
  -v="": Comma separated string of variables to plot
```

##Examples:
In the examples, I am using `ngspice`. You can use `./testobot.sh examples/2b` to automatically generate plots with ngspice.

---

```
ngspice -b examples/2b.net -r examples/2b.raw >/dev/null
./spiceplot -e 'v(b)-v(e)' examples/2b.raw examples/2b.svg
rm examples/2b.raw
```

Output:
```
*hw3.2b: Operating Point
v(vdd) = 2
v(b) = 2
v(e) = 1.2791278613171688
i(vb) = 8.471046155287428e-06
i(v1) = -0.0012791278633171779
v(b)-v(e) = 0.7208721386828312
```

---
```
ngspice -b examples/oscillator.net -r examples/oscillator.raw >/dev/null
./spiceplot examples/oscillator.raw examples/oscillator.svg
rm examples/oscillator.raw
```

Output:
```
*cmos ring oscillator example: Transient Analysis
Outputting  examples/oscillator-0-voltage.svg
Outputting  examples/oscillator-0-current.svg
time = 0
v(vdd) = 3
v(vss) = 0
v(1) = 1.7573593126682514
v(2) = 1.7573593126682516
v(3) = 1.7573593126682516
v(4) = 1.7573593126682516
v(5) = 1.7573593126682516
i(v2) = 0.004411873942866134
i(v1) = -0.0044118739428661395
```

![](http://norcalli.com/spiceplot/oscillator-all-voltage.svg)
![](http://norcalli.com/spiceplot/oscillator-all-current.svg)
---
```
ngspice -b examples/oscillator.net -r examples/oscillator.raw >/dev/null
./spiceplot -v 'i(v1)' examples/oscillator.raw examples/oscillator.svg
rm examples/oscillator.raw
```

Output:
```
*cmos ring oscillator example: Transient Analysis
Plotting 1 variables: Set{i(v1)}
i(v1)
Outputting  examples/oscillator-0-current.svg
time = 0
v(vdd) = 3
v(vss) = 0
v(1) = 1.7573593126682514
v(2) = 1.7573593126682516
v(3) = 1.7573593126682516
v(4) = 1.7573593126682516
v(5) = 1.7573593126682516
i(v2) = 0.004411873942866134
i(v1) = -0.0044118739428661395
```

![](http://norcalli.com/spiceplot/oscillator-variable.svg)
