package main

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"code.google.com/p/plotinum/vg"
	"flag"
	"fmt"
	"github.com/deckarep/golang-set"
	raw "github.com/norcalli/rawspice"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// func MakeXYs(x, y []float64) plotter.XYs {
func SpiceToXY(x, y *raw.SpiceVector) plotter.XYs {
	pts := make(plotter.XYs, len(x.Data))
	for i, _ := range pts {
		pts[i].X, pts[i].Y = x.Data[i], y.Data[i]
	}
	return pts
}

// // func Plot(sp *raw.SpicePlot) []*plot.Plot {
// func Plot(sp *raw.SpicePlot) map[string]*plot.Plot {
// 	// plots := []*plot.Plot{}
// 	plots := map[string]*plot.Plot{}

// 	typemap := map[string][]*raw.SpiceVector{}
// 	// Separate vectors by type into map
// 	for _, v := range sp.Vectors {
// 		typemap[v.Type] = append(typemap[v.Type], v)
// 	}
// 	// If we have a plot which has a time element, then
// 	// we plot it by time.
// 	if timevectors, exists := typemap["time"]; exists {
// 		timevector := timevectors[0]
// 		delete(typemap, "time") // remove time to simplify iteration
// 		// Create different plot for each type.
// 		for vector_type, vectors := range typemap {
// 			plot, err := plot.New()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			for i, vector := range vectors {
// 				xys := SpiceToXY(timevector, vector)
// 				line, err := plotter.NewLine(xys)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				line.Color = plotutil.Color(i + 1)
// 				line.Dashes = plotutil.Dashes(0)
// 				plot.Add(line)
// 				plot.Legend.Add(vector.Name, line)
// 			}
// 			plot.Title.Text = sp.Title + ": " + sp.Name
// 			plot.X.Label.Text = "time"
// 			plot.Y.Label.Text = vector_type
// 			// plots = append(plots, plot)
// 			plots[vector_type] = plot
// 		}
// 	}
// 	return plots
// }

type PlotMap map[string]*plot.Plot

func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func maxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}
func findMultiplicity(vector []float64) int {
	n := len(vector)
	// log.Println("N:", n)
	// l := n
	x0 := vector[0]
	maxdist, maxcount, p := 1, 0, 0
	for i, x := range vector {
		if x == x0 {
			maxcount++
			// maxdist = mathutil.Max(i-p, maxdist)
			maxdist = maxInt(i-p, maxdist)
			p = i
		}
	}

	// log.Println("maxdist:", maxdist)
	// log.Println("maxcount:", maxcount)
	// max := minInt(n/maxdist, maxcount)
	// max := mathutil.Min(n/maxdist, maxcount)
	max := minInt(n/maxdist, maxcount)
	{
		// x0 := vector[0]
		// 	factors := mathutil.FactorInt(uint32(n))
		// loop1:
		// 	for _, f := range factors {
		// 		count := f.Prime
		// 		y := n / count
		// 		for j := 0; j < count; j++ {
		// 			for k := 0; k < y; k++ {
		// 				if vector[j*y+k] != vector[k] {
		// 					continue loop1
		// 				}
		// 			}
		// 		}
		// 		return count
		// 	}
		// 	return 1

	loop2:
		for count := 2; count <= max; count++ {
			if n%count == 0 {
				y := n / count
				for j := 0; j < count; j++ {
					for k := 0; k < y; k++ {
						if vector[j*y+k] != vector[k] {
							continue loop2
						}
					}
				}
				return count
			}
		}
		// log.Println("Ops:", l)
		return 1
	}
}

func plotBehind(sp *raw.SpicePlot, scale_vector *raw.SpiceVector, typemap map[string][]*raw.SpiceVector) PlotMap {
	m := findMultiplicity(scale_vector.Data)
	c := int(sp.NPoints) / m
	plots := PlotMap{}
	for vector_type, vectors := range typemap {
		plot, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		for i, vector := range vectors {
			xys := SpiceToXY(scale_vector, vector)
			for j := 0; j < m; j++ {
				line, err := plotter.NewLine(xys[j*c : (j+1)*c])
				if err != nil {
					log.Fatal(err)
				}
				line.Color = plotutil.Color(i)
				line.Dashes = plotutil.Dashes(0)
				plot.Add(line)
				// plot.Legend.Add(vector.Name, line)
				// plot.Legend.Add(fmt.Sprintf("%s-%d", vector.Name, j), line)
				if j == 0 {
					plot.Legend.Add(vector.Name, line)
				}
			}
			// xys := SpiceToXY(scale_vector, vector)
			// line, err := plotter.NewLine(xys)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// line.Color = plotutil.Color(i + 1)
			// line.Dashes = plotutil.Dashes(0)
			// plot.Add(line)
			// plot.Legend.Add(vector.Name, line)
		}
		plot.Title.Text = sp.Title + ": " + sp.Name
		plot.X.Label.Text = scale_vector.Name
		plot.Y.Label.Text = vector_type
		// plots = append(plots, plot)
		plots[vector_type] = plot
	}
	return plots
}

func Plot(sp *raw.SpicePlot) PlotMap {
	if sp.NPoints == 1 {
		return nil
	}
	data_vectors := sp.Vectors
	scale_vector := data_vectors[0]
	typemap := map[string][]*raw.SpiceVector{}
	// Separate vectors by type into map
	for _, v := range data_vectors[1:] {
		typemap[v.Type] = append(typemap[v.Type], v)
	}
	return plotBehind(sp, scale_vector, typemap)
}

func PlotSome(sp *raw.SpicePlot, vars mapset.Set) PlotMap {
	if sp.NPoints == 1 {
		return nil
	}
	data_vectors := sp.Vectors
	scale_vector := data_vectors[0]

	typemap := map[string][]*raw.SpiceVector{}
	// Separate vectors by type into map
	for _, v := range data_vectors[1:] {
		if vars.Contains(v.Name) {
			fmt.Println(v.Name)
			typemap[v.Type] = append(typemap[v.Type], v)
		}
	}

	return plotBehind(sp, scale_vector, typemap)
}

func usage(name string) {
	fmt.Printf("Usage: %s [flags] <rawfile> <output>\n", name)
	fmt.Println("\nSupported formats for the output are: .svg, .png, .jpg, .jpeg, .eps, .tiff, .pdf")
	fmt.Println("Possible flags:")
	flag.PrintDefaults()
	os.Exit(0)
}

var plotvars, expvars string

func init() {
	flag.StringVar(&plotvars, "v", "", "Comma separated string of variables to evaluate")
	flag.StringVar(&expvars, "e", "", "Comma separated string of expressions to plot")
	// flag.StringVar(&format, "f", "svg", "Output format. Supported formats are: svg,png,pdf,jpg")
}

type VariableTable map[string]float64

func PlotToVariableTable(plot *raw.SpicePlot) (VariableTable, error) {
	// if plot.NPoints != 1 {
	// 	return nil, fmt.Errorf("plot.NPoints > 1")
	// }
	vt := VariableTable{}
	for _, v := range plot.Vectors {
		vt[v.Name] = v.Get(0)
	}
	return vt, nil
}

// This is a hack, I have to eventually come up with dynamic evaluation
// of expressions in Go. I'll probably write a library for it.
func eval(e string, vt VariableTable) (float64, error) {
	for name, value := range vt {
		e = strings.Replace(e, name, fmt.Sprintf("(%g)", value), -1)
	}
	cmd := exec.Command("bc", "-ql")
	cmd.Stdin = strings.NewReader("scale=15;" + e + "\n")
	result, err := cmd.CombinedOutput()
	result = result[:len(result)-1]
	if err != nil {
		return 0, err
	}
	if v, err := strconv.ParseFloat(string(result), 64); err == nil {
		return v, nil
	} else {
		return 0, err
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		usage(os.Args[0])
	}
	filename := flag.Arg(0)
	outputname := flag.Arg(1)
	variables := mapset.NewSet()
	for _, x := range strings.Split(plotvars, ",") {
		if x != "" {
			variables.Add(x)
		}
	}
	expressions := mapset.NewSet()
	for _, x := range strings.Split(expvars, ",") {
		if x != "" {
			expressions.Add(x)
		}
	}
	// variables := mapset.NewSetFromSlice(strings.Split(plotvars, ","))
	// file, err := os.Open(filename, "rb")
	spice_plots, err := raw.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var prefix, extension string
	{
		i := strings.LastIndex(outputname, ".")
		prefix, extension = outputname[:i], outputname[i+1:]
	}
	var outplots map[string]*plot.Plot
	for n, plot := range spice_plots {
		fmt.Printf("%s: %s\n", plot.Title, plot.Name)
		// fmt.Printf("%+v\n", plot)
		// outplots := Plot(plots[len(plots)-1])
		if variables.Cardinality() > 0 {
			fmt.Println("Plotting", variables.Cardinality(), "variables:", variables)
			outplots = PlotSome(plot, variables)
		} else {
			outplots = Plot(plot)
		}
		for i, plot := range outplots {
			newname := fmt.Sprintf("%s-%d-%v.%s", prefix, n, i, extension)
			mult := 1.0
			multvg := vg.Length(mult)
			initLineWidth := plot.X.LineStyle.Width
			initFontSize := plot.X.Label.Font.Size
			initTickFontSize := plot.X.Tick.Label.Font.Size
			initTickLength := plot.X.Tick.Width
			plot.X.Label.Font.Size = multvg * initFontSize
			plot.Y.Label.Font.Size = multvg * initFontSize
			plot.X.Tick.Label.Font.Size = multvg * initTickFontSize
			plot.Y.Tick.Label.Font.Size = multvg * initTickFontSize
			plot.X.Tick.Width = multvg * initTickLength
			plot.Y.Tick.Width = multvg * initTickLength
			plot.X.LineStyle.Width = multvg * initLineWidth
			plot.Y.LineStyle.Width = multvg * initLineWidth

			fmt.Println("Outputting ", newname)
			if err := plot.Save(6*mult, 4*mult, newname); err != nil {
				panic(err)
			}
		}
		for _, x := range plot.Vectors {
			fmt.Printf("%s = %v\n", x.Name, x.Get(0))
		}
		if expressions.Cardinality() > 0 { // && plot.NPoints == 1 {
			vt, err := PlotToVariableTable(plot)
			if err != nil {
				log.Fatal(err)
			}
			for x := range expressions.Iter() {
				s := x.(string)
				// If there is an error, just continue.
				if r, err := eval(s, vt); err == nil {
					fmt.Printf("%s = %v\n", s, r)
				}
			}
		}
		fmt.Println()
	}
}
