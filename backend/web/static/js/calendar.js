// Custom calendar renderer for trip schedules
var currentDate = (function() {
    var s = window.tripSchedules;
    if (s && s.length > 0 && s[0].DateStart) {
        var d = new Date(s[0].DateStart + 'T00:00:00');
        if (!isNaN(d.getTime())) return d;
    }
    return new Date();
})();
var schedules = window.tripSchedules || [];

function changeMonth(delta) {
    currentDate.setMonth(currentDate.getMonth() + delta);
    renderCalendar();
    highlightAllSchedules();
}

function renderCalendar() {
    var year = currentDate.getFullYear();
    var month = currentDate.getMonth();

    var monthNames = ['Januari','Februari','Maret','April','Mei','Juni','Juli','Agustus','September','Oktober','November','Desember'];
    document.getElementById('calendarMonth').textContent = monthNames[month] + ' ' + year;

    var firstDay = new Date(year, month, 1);
    var startDate = new Date(firstDay);
    startDate.setDate(startDate.getDate() - ((firstDay.getDay() + 6) % 7));

    var today = new Date();
    var todayStr = today.getFullYear() + '-' + String(today.getMonth()+1).padStart(2,'0') + '-' + String(today.getDate()).padStart(2,'0');

    var grid = document.getElementById('calendarGrid');
    grid.innerHTML = '';

    for (var w = 0; w < 6; w++) {
        var weekDiv = document.createElement('div');
        weekDiv.className = 'week-row';
        weekDiv.style.display = 'grid';
        weekDiv.style.gridTemplateColumns = 'repeat(7, 1fr)';

        var weekEvents = [];
        var weekStart = new Date(startDate);
        weekStart.setDate(weekStart.getDate() + w * 7);
        var weekEnd = new Date(weekStart);
        weekEnd.setDate(weekEnd.getDate() + 6);

        for (var d = 0; d < 7; d++) {
            var cellDate = new Date(weekStart);
            cellDate.setDate(cellDate.getDate() + d);
            var dateStr = cellDate.getFullYear() + '-' + String(cellDate.getMonth()+1).padStart(2,'0') + '-' + String(cellDate.getDate()).padStart(2,'0');

            var dayCell = document.createElement('div');
            dayCell.className = 'day-cell';
            if (cellDate.getMonth() !== month) dayCell.classList.add('other-month');
            if (dateStr === todayStr) dayCell.classList.add('today');

            dayCell.innerHTML = '<div class="day-number">' + cellDate.getDate() + '</div>';
            dayCell.dataset.date = dateStr;

            // Click handler for day
            dayCell.addEventListener('click', function(e) {
                if (e.target.closest('.event-bar')) return;
                onDayClick(dateStr);
            });

            weekDiv.appendChild(dayCell);

            // Check event overlap for this day
            schedules.forEach(function(s) {
                var sStart = s.DateStart;
                var sEnd = s.DateEnd;
                if (dateStr >= sStart && dateStr <= sEnd) {
                    var evIdx = -1;
                    for (var ei = 0; ei < weekEvents.length; ei++) {
                        if (weekEvents[ei].id === s.ID) {
                            evIdx = ei;
                            break;
                        }
                    }
                    if (evIdx === -1) {
                        var evStart = sStart < weekStart.toISOString().slice(0,10) ? weekStart.toISOString().slice(0,10) : sStart;
                        var evEnd = sEnd > weekEnd.toISOString().slice(0,10) ? weekEnd.toISOString().slice(0,10) : sEnd;
                        var startIdx = 0;
                        var endIdx = 6;
                        var tmp = new Date(weekStart);
                        for (var i = 0; i < 7; i++) {
                            var ds = tmp.getFullYear() + '-' + String(tmp.getMonth()+1).padStart(2,'0') + '-' + String(tmp.getDate()).padStart(2,'0');
                            if (ds === evStart) startIdx = i;
                            if (ds === evEnd) endIdx = i;
                            tmp.setDate(tmp.getDate() + 1);
                        }
                        weekEvents.push({
                            id: s.ID,
                            startIndex: startIdx,
                            endIndex: endIdx,
                            dateStart: s.DateStart,
                            dateEnd: s.DateEnd
                        });
                    }
                }
            });
        }

        // Row assignment for overlapping events
        if (weekEvents.length > 0) {
            weekEvents.sort(function(a, b) { return a.startIndex - b.startIndex || (b.endIndex - b.startIndex) - (a.endIndex - a.startIndex); });

            var rowAssignments = [];
            weekEvents.forEach(function(ev) {
                var row = 0;
                var placed = false;
                while (!placed) {
                    var conflict = false;
                    if (rowAssignments[row]) {
                        for (var ri = 0; ri < rowAssignments[row].length; ri++) {
                            var assigned = rowAssignments[row][ri];
                            if (!(ev.endIndex < assigned.startIndex || ev.startIndex > assigned.endIndex)) {
                                conflict = true;
                                break;
                            }
                        }
                    }
                    if (!conflict) {
                        if (!rowAssignments[row]) rowAssignments[row] = [];
                        rowAssignments[row].push(ev);
                        ev.row = row;
                        placed = true;
                    } else {
                        row++;
                    }
                }
            });

            weekEvents.forEach(function(ev) {
                var span = ev.endIndex - ev.startIndex + 1;
                var bar = document.createElement('div');
                var colorClass = 'event-color-' + ((ev.id % 5) + 1);
                var barClass = 'event-bar ' + colorClass + ' row-' + ev.row;
                if (span === 1) {
                    barClass += ' event-bar-single';
                } else {
                    barClass += ' event-bar-start';
                }

                bar.className = barClass;
                bar.style.left = (ev.startIndex * (100/7)) + '%';
                bar.style.width = (span * (100/7)) + '%';
                bar.title = ev.dateStart + ' - ' + ev.dateEnd;
                bar.textContent = formatDateShort(ev.dateStart) + (ev.dateStart !== ev.dateEnd ? '-' + formatDateShort(ev.dateEnd) : '');

                bar.addEventListener('click', function() {
                    window.location.href = '/schedules/' + ev.id + '/edit';
                });

                weekDiv.appendChild(bar);
            });
        }

        grid.appendChild(weekDiv);
        if (w === 0) {
            startDate.setDate(startDate.getDate() + 7);
        }
    }
}

function formatDateShort(dateStr) {
    var parts = dateStr.split('-');
    var months = ['','Jan','Feb','Mar','Apr','Mei','Jun','Jul','Agu','Sep','Okt','Nov','Des'];
    var m = parseInt(parts[1], 10);
    return parseInt(parts[2], 10) + ' ' + (months[m] || '');
}

// Init calendar on page load
document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById('calendarGrid')) {
        renderCalendar();
        highlightAllSchedules();
    }
});

// ── Day click ──

function onDayClick(dateStr) {
    var matching = schedules.filter(function(s) {
        return dateStr >= s.DateStart && dateStr <= s.DateEnd;
    });
    if (matching.length === 0) return;

    var info = document.getElementById('scheduleInfo');
    if (matching.length === 1) {
        window.location.href = '/schedules/' + matching[0].ID + '/edit';
        return;
    }

    info.textContent = 'Menampilkan ' + matching.length + ' jadwal pada ' + formatDateShort(dateStr);
    document.querySelectorAll('.sched-item').forEach(function(item) {
        var start = item.dataset.start;
        var end = item.dataset.end;
        var match = dateStr >= start && dateStr <= end;
        item.style.display = match ? '' : 'none';
        if (match) {
            item.classList.add('ring-2', 'ring-teal-400');
        } else {
            item.classList.remove('ring-2', 'ring-teal-400');
        }
    });
}

function highlightAllSchedules() {
    var info = document.getElementById('scheduleInfo');
    if (info) info.textContent = 'Menampilkan semua jadwal';
    document.querySelectorAll('.sched-item').forEach(function(item) {
        item.style.display = '';
        item.classList.remove('ring-2', 'ring-teal-400');
    });
}
