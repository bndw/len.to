import {data} from '@params';
if (!data) {
  return
}

const calData = Object.keys(data).map(date => ({
    date: date,
    total: data[date],
    details: [],
    summary: []
}))

const div_id = 'calendar';
const label = 'Photos';
const color = '#cd2327';
const overview = 'global'; // global, year, month, day
const handler = function (val) {
  // TODO: What does this do?
  console.log(val);
  window.location = "/about"
};

// Initialize calendar heatmap
calendarHeatmap.init(calData, div_id, color, overview, handler, label);
