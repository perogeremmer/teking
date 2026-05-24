// Dashboard chart
document.addEventListener('DOMContentLoaded', function() {
    var canvas = document.getElementById('dashboardChart');
    if (!canvas) return;

    var monthsData = window.dashboardData || null;
    if (!monthsData) return;

    var monthNames = ['Jan','Feb','Mar','Apr','Mei','Jun','Jul','Agu','Sep','Okt','Nov','Des'];
    var labels = monthsData.map(function(d) { return monthNames[parseInt(d.Month)-1] || d.Month; });

    new Chart(canvas.getContext('2d'), {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Revenue',
                    data: monthsData.map(function(d) { return d.Revenue; }),
                    backgroundColor: '#0d9488',
                    borderRadius: 4,
                    yAxisID: 'y'
                },
                {
                    label: 'Bookings',
                    data: monthsData.map(function(d) { return d.Bookings; }),
                    backgroundColor: '#f97316',
                    borderRadius: 4,
                    yAxisID: 'y1'
                }
            ]
        },
        options: {
            responsive: true,
            interaction: {
                mode: 'index',
                intersect: false
            },
            plugins: {
                legend: { position: 'top' }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    position: 'left',
                    ticks: {
                        callback: function(v) {
                            if (v >= 1000000) return (v/1000000).toFixed(1)+'jt';
                            if (v >= 1000) return (v/1000).toFixed(0)+'rb';
                            return v;
                        }
                    }
                },
                y1: {
                    beginAtZero: true,
                    position: 'right',
                    grid: { drawOnChartArea: false },
                    ticks: {
                        stepSize: 1
                    }
                }
            }
        }
    });
});
