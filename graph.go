package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"text/template"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa/ssautil"

	"github.com/gophergala2016/traceapp/cha"
	"github.com/gophergala2016/traceapp/node"
	"github.com/gophergala2016/traceapp/rta"
	"github.com/julienschmidt/httprouter"
)

func Grapher(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	repo := ps.ByName("importpath")
	repo = strings.TrimLeft(repo, "/")
	log.Println(repo)
	_, err := exec.Command("go", "get", "-u", "-f", repo).Output()
	if err != nil {
		w.Write([]byte("couldn't get repo [" + repo + "]:" + err.Error()))
		return
	}

	var conf loader.Config
	conf.Import(repo)
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	ssaProg := ssautil.CreateProgram(prog, 0)
	ssaProg.Build()

	var nodes []node.Node
	nodes, err = rta.GetNodes(ssaProg)
	if err != nil {
		nodes, err = cha.GetNodes(ssaProg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	by, err := json.Marshal(nodes)
	if err != nil {
		panic(err)
	}

	var data = struct {
		Graph string
	}{
		Graph: string(by),
	}

	tmpl, _ := template.New("foo").Parse(nt)
	tmpl.Execute(w, data)
}

var nt = `
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8"/>
    <link type="text/css" rel="stylesheet" href="/static/css/style.css"/>
    <style type="text/css">
path.arc {
  cursor: move;
  fill: #fff;
}
.node {
  font-size: 10px;
}
.node:hover {
  fill: #1f77b4;
}
.link {
  fill: none;
  stroke: #1f77b4;
  stroke-opacity: .4;
  pointer-events: none;
}
.link.source, .link.target {
  stroke-opacity: 1;
  stroke-width: 2px;
}
.node.target {
  fill: #d62728 !important;
}
.link.source {
  stroke: #d62728;
}
.node.source {
  fill: #2ca02c;
}
.link.target {
  stroke: #2ca02c;
}
    </style>
  </head>
  <body>
    <div style="position:absolute;bottom:0;font-size:18px;">tension: <input style="position:relative;top:3px;" type="range" min="0" max="100" value="85"></div>
    <script type="text/javascript" src="/static/js/d3.js"></script>
    <script type="text/javascript" src="/static/js/d3.layout.js"></script>
    <script type="text/javascript" src="/static/js/packages.js"></script>
    <script type="text/javascript">
	var win = window,
	    doc = document,
	    e = doc.documentElement,
	    g = doc.getElementsByTagName('body')[0],
	    x = win.innerWidth || e.clientWidth || g.clientWidth,
	    y = win.innerHeight|| e.clientHeight|| g.clientHeight;
var w = x / (5/4),
    h = y / (5/4),
    rx = w / 2,
    ry = h / 2,
    m0,
    rotate = 0;
var splines = [];
var cluster = d3.layout.cluster()
    .size([360, ry - 120])
    .sort(function(a, b) { return d3.ascending(a.key, b.key); });
var bundle = d3.layout.bundle();
var line = d3.svg.line.radial()
    .interpolate("bundle")
    .tension(.85)
    .radius(function(d) { return d.y; })
    .angle(function(d) { return d.x / 180 * Math.PI; });
// Chrome 15 bug: <http://code.google.com/p/chromium/issues/detail?id=98951>
var div = d3.select("body").insert("div", "h2")
    .style("top", "-80px")
    .style("left", "-160px")
    .style("width", w + "px")
    .style("height", w + "px")
    .style("position", "absolute")
    .style("-webkit-backface-visibility", "hidden");
var svg = div.append("svg:svg")
    .attr("width", w)
    .attr("height", w)
  .append("svg:g")
    .attr("transform", "translate(" + rx + "," + ry + ")");
svg.append("svg:path")
    .attr("class", "arc")
    .attr("d", d3.svg.arc().outerRadius(ry - 120).innerRadius(0).startAngle(0).endAngle(2 * Math.PI))
    .on("mousedown", mousedown);
classes = JSON.parse('{{ .Graph }}')
  var nodes = cluster.nodes(packages.root(classes)),
      links = packages.imports(nodes),
      splines = bundle(links);
  var path = svg.selectAll("path.link")
      .data(links)
    .enter().append("svg:path")
      .attr("class", function(d) { return "link source-" + d.source.key + " target-" + d.target.key; })
      .attr("d", function(d, i) { return line(splines[i]); });
  svg.selectAll("g.node")
      .data(nodes.filter(function(n) { return !n.children; }))
    .enter().append("svg:g")
      .attr("class", "node")
      .attr("id", function(d) { return "node-" + d.key; })
      .attr("transform", function(d) { return "rotate(" + (d.x - 90) + ")translate(" + d.y + ")"; })
    .append("svg:text")
      .attr("dx", function(d) { return d.x < 180 ? 8 : -8; })
      .attr("dy", ".31em")
      .attr("text-anchor", function(d) { return d.x < 180 ? "start" : "end"; })
      .attr("transform", function(d) { return d.x < 180 ? null : "rotate(180)"; })
      .text(function(d) { return d.key; })
      .on("mouseover", mouseover)
      .on("mouseout", mouseout);
  d3.select("input[type=range]").on("change", function() {
    line.tension(this.value / 100);
    path.attr("d", function(d, i) { return line(splines[i]); });
  });
d3.select(window)
    .on("mousemove", mousemove)
    .on("mouseup", mouseup);
function mouse(e) {
  return [e.pageX - rx, e.pageY - ry];
}
function mousedown() {
  m0 = mouse(d3.event);
  d3.event.preventDefault();
}
function mousemove() {
  if (m0) {
    var m1 = mouse(d3.event),
        dm = Math.atan2(cross(m0, m1), dot(m0, m1)) * 180 / Math.PI;
    div.style("-webkit-transform", "translateY(" + (ry - rx) + "px)rotateZ(" + dm + "deg)translateY(" + (rx - ry) + "px)");
  }
}
function mouseup() {
  if (m0) {
    var m1 = mouse(d3.event),
        dm = Math.atan2(cross(m0, m1), dot(m0, m1)) * 180 / Math.PI;
    rotate += dm;
    if (rotate > 360) rotate -= 360;
    else if (rotate < 0) rotate += 360;
    m0 = null;
    div.style("-webkit-transform", null);
    svg
        .attr("transform", "translate(" + rx + "," + ry + ")rotate(" + rotate + ")")
      .selectAll("g.node text")
        .attr("dx", function(d) { return (d.x + rotate) % 360 < 180 ? 8 : -8; })
        .attr("text-anchor", function(d) { return (d.x + rotate) % 360 < 180 ? "start" : "end"; })
        .attr("transform", function(d) { return (d.x + rotate) % 360 < 180 ? null : "rotate(180)"; });
  }
}
function mouseover(d) {
  svg.selectAll("path.link.target-" + d.key)
      .classed("target", true)
      .each(updateNodes("source", true));
  svg.selectAll("path.link.source-" + d.key)
      .classed("source", true)
      .each(updateNodes("target", true));
}
function mouseout(d) {
  svg.selectAll("path.link.source-" + d.key)
      .classed("source", false)
      .each(updateNodes("target", false));
  svg.selectAll("path.link.target-" + d.key)
      .classed("target", false)
      .each(updateNodes("source", false));
}
function updateNodes(name, value) {
  return function(d) {
    if (value) this.parentNode.appendChild(this);
    svg.select("#node-" + d[name].key).classed(name, value);
  };
}
function cross(a, b) {
  return a[0] * b[1] - a[1] * b[0];
}
function dot(a, b) {
  return a[0] * b[0] + a[1] * b[1];
}
    </script>
  </body>
</html>
`
