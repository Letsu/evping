<template>
  <div class="graph-container">
    <Line
      id="my-chart-id"
      v-if="loaded"
      :options="chartOptions"
      :data="chartData"
    />
    <div class="d-flex statistics">
      <p class="mr-2 mt-2">Min: {{ statistics.min }} ms (at {{ statistics.minTime }})</p>
      <p class="ma-2">Max: {{ statistics.max }} ms (at {{ statistics.maxTime }})</p>
      <p class="ma-2">Mean: {{ statistics.mean }} ms</p>
      <p class="ma-2">Lost Packages: {{ statistics.lostPercentage }}%</p>
    </div>
  </div>
</template>

<script>
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement } from 'chart.js'
import axios from 'axios'
import moment from 'moment';
import {CrosshairPlugin,Interpolate} from 'chartjs-plugin-crosshair';

ChartJS.register(Title, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement)
ChartJS.register(CrosshairPlugin);


function down(ctx, value) {
  return ctx.p0.skip || ctx.p1.skip || ctx.p1.parsed.y === -1 ? 'red' : undefined;
}

export default {
  name: 'BarChart',
  components: { Line },
  data() {
    return {
      loaded: false,
      intervalId: null,
      statistics: {
        min: null,
        max: null,
        mean: null,
        lostPercentage: null,
        minTime: null,
        maxTime: null,
      },
      chartData: {
        labels: [],
        datasets: [ 
          { 
            label: '',
            backgroundColor: '#f87979',
            borderColor: 'lightgreen',
            tension: 0.1,
            fill: false,
            spanGaps: true,
            pointRadius: 0,
            data: [],
            stepped: true,
            segment: {
              borderColor: ctx => down(ctx, 'rgb(192,75,75)'),
            }
          } 
        ]
      },
      chartOptions: {
        animation: false,
        scales: {
          x: {
            ticks: {
              autoSkip: true,
              maxTicksLimit: 10 // adjust this to the number of labels you want to show
            }
          },
          y: {
            min: -2,
          }
        },
        maintainAspectRatio: true,
        responsive: true,
        fill: false,
        interaction: {
          intersect: false
        },
        plugins: {
          tooltip: {
            mode: 'interpolate',
            intersect: false
          },
          crosshair: {
            line: {
              color: '#F66',  // crosshair line color
              width: 1        // crosshair line width
            },
            sync: {
              enabled: true,            // enable trace line syncing with other charts
              group: 1,                 // chart group
              suppressTooltips: false   // suppress tooltips when showing a synced tracer
            },
            zoom: {
              enabled: true,                                      // enable zooming
              zoomboxBackgroundColor: 'rgba(66,133,244,0.2)',     // background color of zoom box 
              zoomboxBorderColor: '#48F',                         // border color of zoom box
              zoomButtonText: 'Reset Zoom',                       // reset zoom button text
              zoomButtonClass: 'reset-zoom',                      // reset zoom button class
            },
            callbacks: {
              beforeZoom: this.beforeZoom,
            }
          }
        },
      },
    }
  },
  methods: {
    formatTime(time) {
      let date = new Date(time);
      let day = date.getDate().toString().padStart(2, '0');
      let month = (date.getMonth() + 1).toString().padStart(2, '0'); // Months are 0-based
      let hours = date.getHours().toString().padStart(2, '0');
      let minutes = date.getMinutes().toString().padStart(2, '0');
      let seconds = date.getSeconds().toString().padStart(2, '0');

      return `${day}.${month} ${hours}:${minutes}:${seconds}`;
    },
    async loadData() {
      try {
        const response = await axios.get('http://127.0.0.1:8080/api/');
        console.log('Response data:', response.data);
        if (response.data[0].labels && response.data[0].data) {
          //this.chartData.labels = response.data[0].labels.map(label => {
          //  let date = new Date(label);
          //  return date.toLocaleString();
          //});
          let labels = response.data[0].labels;
          let labels_time = labels.map(label => {
            let date = moment(label);
            return date.fromNow();
          });
          let data = response.data[0].data
          let label = response.data[0].host;

          // Calculate statistics
          let filteredData = data.filter(value => value !== -1);

          let min = Math.min(...filteredData);
          let max = Math.max(...filteredData);
          let mean = filteredData.reduce((a, b) => a + b, 0) / filteredData.length;

          let minIndex = data.indexOf(min);
          let maxIndex = data.indexOf(max);

          let minTime = this.formatTime(labels[minIndex]);
          let maxTime = this.formatTime(labels[maxIndex]);

          let lostPackages = data.filter(value => value === -1).length;
          let totalPackages = data.length;
          let lostPercentage = (lostPackages / totalPackages) * 100;


          //let data = response.data[0].data.map(value => value === -1 ? NaN : value);
          this.chartData = {
            labels: labels_time,
            datasets: [ 
              { 
                label: label,
                backgroundColor: '#f87979',
                borderColor: 'lightgreen',
                tension: 0.1,
                fill: false,
                spanGaps: true,
                pointRadius: 0,
                data: data,
                stepped: true,
                segment: {
                  borderColor: ctx => down(ctx, 'rgb(192,75,75)'),
                }
              } 
            ]
          }
          this.statistics = {
            min: min,
            max: max,
            mean: Number(mean).toFixed(1),
            lostPercentage: Number(lostPercentage).toFixed(2),
            minTime: minTime,
            maxTime: maxTime,
          }
          this.loaded = true;
        } else {
          console.error('Unexpected response structure:', response.data);
        }
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    },
    startInterval() {
      this.intervalId = setInterval(this.loadData, 100000);
      //clearInterval(this.intervalId);
    },
    stopInterval() {
      clearInterval(this.intervalId);
      this.intervalId = null;
    },
    beforeZoom: function() {
      const self = this;
      return function(start, end) {  
        self.stopInterval();
        self.$nextTick(() => {
          const resetZoomButton = document.querySelector('.reset-zoom');
          if (resetZoomButton) {
            resetZoomButton.addEventListener('click', self.startInterval);
          }
        });
        return true;
      };
    },
  },
  mounted() {
    this.loaded = false
    this.loadData();
    this.startInterval();
  },
  beforeDestroy() {
    this.stopInterval();
  },
}
</script>

<style scoped>
  button {
    background-color: #ff3784;
    border: none;
    border-radius: 0;
    color: white;
    cursor: pointer;
    padding: 0.4rem 1rem;
  }

  button:hover {
    background-color: #f0287a;
    color: white;
    text-decoration: none;
  }

  .graph-container {
    position: relative;
  }

  .statistics {
    position: absolute;
    bottom: 5%;
    left: 5%;
  }
</style>