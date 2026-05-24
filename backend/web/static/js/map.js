// Leaflet map for meeting point form
document.addEventListener('DOMContentLoaded', function() {
    var mapEl = document.getElementById('map');
    if (!mapEl) return;

    var map = L.map('map').setView([window.mpLat || -6.8947, window.mpLng || 108.4], 10);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    var marker = L.marker([window.mpLat, window.mpLng], { draggable: true }).addTo(map);

    marker.on('dragend', function() {
        var latlng = marker.getLatLng();
        document.getElementById('lat').value = latlng.lat.toFixed(6);
        document.getElementById('lng').value = latlng.lng.toFixed(6);
    });

    map.on('click', function(e) {
        marker.setLatLng(e.latlng);
        document.getElementById('lat').value = e.latlng.lat.toFixed(6);
        document.getElementById('lng').value = e.latlng.lng.toFixed(6);
    });

    // Geocoding search
    var searchEl = document.createElement('div');
    searchEl.style.cssText = 'position: absolute; top: 10px; left: 10px; z-index: 1000; width: calc(100% - 20px);';
    searchEl.innerHTML = '<input type="text" id="mapSearch" placeholder="Cari lokasi..." class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg shadow">';
    mapEl.parentElement.style.position = 'relative';
    mapEl.parentElement.appendChild(searchEl);

    document.getElementById('mapSearch').addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            var q = this.value;
            if (!q) return;
            fetch('https://nominatim.openstreetmap.org/search?format=json&q=' + encodeURIComponent(q) + '&limit=1')
                .then(function(r) { return r.json(); })
                .then(function(data) {
                    if (data && data.length > 0) {
                        var lat = parseFloat(data[0].lat);
                        var lng = parseFloat(data[0].lon);
                        map.setView([lat, lng], 15);
                        marker.setLatLng([lat, lng]);
                        document.getElementById('lat').value = lat.toFixed(6);
                        document.getElementById('lng').value = lng.toFixed(6);
                    }
                })
                .catch(function() {});
        }
    });
});
