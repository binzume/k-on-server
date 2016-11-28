"use strict";
google.charts.load('current', {'packages':['line'], 'language': 'ja'});
google.charts.setOnLoadCallback(drawChart);

var hours = 24;
var devName = "_default";
if (location.search.match(/^\?\w+$/)) {
  devName = location.search.substr(1);
}

function drawChart() {
  $.ajax({
      url: "stats/"+devName+"/values?limit=" + (hours * 12),
      dataType: "json",
      cache : false
  }).done(function(result){
    var rows = [];
    for (var i=0;i <result.values.length; i++) {
      var d = result.values[i];
      rows.push([new Date(d['_timestamp']), d['temp'], d['humid']]);
    }
    var data = new google.visualization.DataTable();
    data.addColumn('date', 'Time');
    data.addColumn('number', 'Temperature');
    data.addColumn('number', 'Humidity');
    data.addRows(rows);

      var options = {
        chart: {
          title: 'k-on : ' + result.description
        },
        height: 500,
        series: {
          // Gives each series an axis name that matches the Y-axis below.
          0: {axis: 'Temp'},
          1: {axis: 'Humid'}
        },
        axes: {
          // Adds labels to each axis; they don't have to match the axis names.
          y: {
            Temp: {label: 'Temp(Celsius)', range:{min:10, max:35}, format:{pattern:'##.##'}},
            Humid: {label: 'Humidity(%)',range:{min:0, max:100}, format:{pattern:'##.#\'%\''}}
          }
        },
        legend: { position: 'bottom' }
      };

    var chart = new google.charts.Line(document.getElementById('chart'));
    chart.draw(data, options);
  });
}

$(document).ready(function(){
  $('#select_term').change(function(){
    hours = $('#select_term').val() | 0;
    drawChart();
  });

  $.ajax({
      url: "stats",
      dataType: "json",
      cache : false
  }).done(function(result){
    var stats = result.stats;
    if (stats == null || stats.length == 0) return;
    for (var i=0; i<stats.length; i++) {
      if (stats[i].display_order < 0) continue;
      $("#menu").append("<nav class='mdl-navigation'><a class='mdl-navigation__link' href='?"+stats[i].name+"'>"+stats[i].description+"</a></nav>");
    }
  });
});
