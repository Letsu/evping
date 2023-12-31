<template>
  <Line
    id="my-chart-id"
    v-if="loaded"
    :options="chartOptions"
    :data="chartData"
  />
</template>

<script>
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement } from 'chart.js'
import axios from 'axios'
import moment from 'moment';
import zoomPlugin from 'chartjs-plugin-zoom';

ChartJS.register(Title, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement)
ChartJS.register(zoomPlugin);


function down(ctx, value) {
  return ctx.p0.skip || ctx.p1.skip || ctx.p1.parsed.y === -1 ? 'red' : undefined;
}

export default {
  name: 'BarChart',
  components: { Line },
  data() {
    return {
      loaded: false,
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
        maintainAspectRatio: false,
        responsive: true,
        fill: false,
        interaction: {
          intersect: false
        },
        plugins: {
          zoom: {
            pan: {
              enabled: true,
              mode: 'x',
            },
            zoom: {
              wheel: {
                enabled: true,
              },
              pinch: {
                enabled: true,
              },
              mode: 'x',
            },
          },
        },
      },
    }
  },
  async mounted() {
    this.loaded = false
    try {
      const response = await axios.get('http://127.0.0.1:8080/api/');
      console.log('Response data:', response.data);
      if (response.data[0].labels && response.data[0].data) {
        //this.chartData.labels = response.data[0].labels.map(label => {
        //  let date = new Date(label);
        //  return date.toLocaleString();
        //});
        this.chartData.labels = response.data[0].labels.map(label => {
          let date = moment(label);
          return date.fromNow();
        });
        let data = response.data[0].data
        //let data = response.data[0].data.map(value => value === -1 ? NaN : value);
        this.chartData.datasets[0].data = data;
        this.chartData.datasets[0].label = response.data[0].host;
        this.loaded = true;
      } else {
        console.error('Unexpected response structure:', response.data);
      }
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  }
}
</script>
