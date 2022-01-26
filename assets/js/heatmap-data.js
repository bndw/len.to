// Copyright 2022 Ben Woodward. All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.
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
