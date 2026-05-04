// Admin Operator Dashboard Logic
// Operator: Tiga Dewa Adventure

const CURRENT_OPERATOR_ID = 'tigadewa';

// Initialize
function initAdmin() {
  loadOperatorData();
  renderSidebar();
  setupMobileMenu();
}

// Load operator data
function loadOperatorData() {
  const operator = getOperatorById(CURRENT_OPERATOR_ID);
  if (operator) {
    document.querySelectorAll('.operator-name').forEach(el => {
      el.textContent = operator.name;
    });
    document.querySelectorAll('.operator-logo').forEach(el => {
      el.src = operator.logo;
    });
  }
}

// Render sidebar
function renderSidebar() {
  const sidebar = document.querySelector('.sidebar-menu');
  if (!sidebar) return;
  
  const currentPage = window.location.pathname.split('/').pop() || 'index.html';
  
  const menuItems = [
    { href: '../index.html', icon: '📊', label: 'Dashboard' },
    { href: '../trips/index.html', icon: '🏔️', label: 'Manajemen Trip' },
    { href: '../schedules/index.html', icon: '📅', label: 'Jadwal Trip' },
    { href: '../bookings/index.html', icon: '📝', label: 'Manajemen Booking' },
    { href: '../reports/revenue.html', icon: '💰', label: 'Laporan' },
    { href: '../profile/index.html', icon: '⚙️', label: 'Pengaturan' },
  ];
  
  sidebar.innerHTML = menuItems.map(item => {
    const isActive = currentPage === item.href.split('/').pop();
    return `
      <a href="${item.href}" class="sidebar-item ${isActive ? 'active' : ''}">
        <span class="icon">${item.icon}</span>
        <span>${item.label}</span>
      </a>
    `;
  }).join('');
}

// Mobile menu toggle
function setupMobileMenu() {
  const toggle = document.querySelector('.menu-toggle');
  const sidebar = document.querySelector('.admin-sidebar');
  const overlay = document.querySelector('.sidebar-overlay');
  
  if (toggle) {
    toggle.addEventListener('click', () => {
      sidebar.classList.toggle('open');
      overlay.classList.toggle('show');
    });
  }
  
  if (overlay) {
    overlay.addEventListener('click', () => {
      sidebar.classList.remove('open');
      overlay.classList.remove('show');
    });
  }
}

// Get operator trips
function getOperatorTrips() {
  return TRIPS.filter(trip => trip.operator_id === CURRENT_OPERATOR_ID);
}

// Get operator bookings
function getOperatorBookings() {
  const operatorTripIds = getOperatorTrips().map(t => t.id);
  return BOOKINGS.filter(booking => operatorTripIds.includes(booking.trip_id));
}

// Calculate revenue
function calculateRevenue(bookings) {
  return bookings.reduce((total, booking) => {
    if (booking.status !== 'cancelled') {
      return total + booking.total;
    }
    return total;
  }, 0);
}

// Format currency
function formatCurrency(amount) {
  return 'Rp ' + amount.toLocaleString('id-ID');
}

// Get status badge HTML
function getStatusBadge(status) {
  const badges = {
    pending: '<span class="badge badge-pending">Pending</span>',
    confirmed: '<span class="badge badge-confirmed">Confirmed</span>',
    completed: '<span class="badge badge-completed">Completed</span>',
    cancelled: '<span class="badge badge-cancelled">Cancelled</span>'
  };
  return badges[status] || status;
}

// Logout
function logout() {
  window.location.href = 'login.html';
}

// Initialize when DOM ready
document.addEventListener('DOMContentLoaded', initAdmin);