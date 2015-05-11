window.onload = function() {
  var width = 1200;
  var chart, svg;
  var data = [];
  function poll() {
    function reqListener () {
      console.log(this.responseText);
      data = JSON.parse(this.responseText);
      d3.select("#timeline svg").remove();
      draw();
    }

    var oReq = new XMLHttpRequest();
    oReq.onload = reqListener;
    oReq.open("get", "/api", true);
    oReq.send();
  }

  function update() {
    svg.selectAll("rect")
      .data([data])
      .transition();
  }

  function draw() {
    chart = d3.timeline()
      .stack()
      .tickFormat({format: d3.time.format("%H:%M:%S"),
            tickTime: d3.time.seconds,
            tickInterval: 10,
            tickSize: 10});

    svg = d3.select("#timeline").append("svg").attr("width", width)
      .data([data])
      .call(chart);
  }

  poll();
  window.setInterval(poll, 50000);
}
