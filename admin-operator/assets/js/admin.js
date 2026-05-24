// Admin Operator Dashboard Logic
// Operator: Tiga Dewa Adventure

const CURRENT_OPERATOR_ID = 'tigadewa';

// Initialize
function initAdmin() {
  loadOperatorData();
  renderSidebar();
  renderSidebarUser();
  setupMobileMenu();
  setupSidebarSearch();
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

// Get base path based on current location depth
function getBasePath() {
  const path = window.location.pathname;
  const adminIndex = path.indexOf('/admin-operator/');
  if (adminIndex === -1) return './';
  
  const subPath = path.substring(adminIndex + '/admin-operator/'.length);
  const depth = subPath.split('/').length - 1;
  
  return depth <= 0 ? './' : '../'.repeat(depth);
}

// Render sidebar user profile card
function renderSidebarUser() {
  const sidebar = document.querySelector('.admin-sidebar');
  if (!sidebar) return;
  
  // Check if user card already exists
  if (sidebar.querySelector('.sidebar-user')) return;
  
  const operator = getOperatorById(CURRENT_OPERATOR_ID);
  const userHTML = `
    <div class="sidebar-user">
      <div class="sidebar-user-card" onclick="logout()">
        <img src="${operator?.logo || ''}" alt="" class="operator-logo">
        <div class="sidebar-user-info">
          <div class="sidebar-user-name operator-name">${operator?.name || 'Operator'}</div>
          <div class="sidebar-user-role">Operator Open Trip</div>
        </div>
        <span class="sidebar-user-arrow"><i class='bx bx-chevron-down'></i></span>
      </div>
    </div>
  `;
  
  sidebar.insertAdjacentHTML('beforeend', userHTML);
}

// Render sidebar
function renderSidebar() {
  const sidebar = document.querySelector('.sidebar-menu');
  if (!sidebar) return;
  
  const currentPage = window.location.pathname.split('/').pop() || 'index.html';
  const basePath = getBasePath();
  
  const menuSections = [
    {
      label: 'Menu Utama',
      items: [
        { href: basePath + 'index.html', icon: "<i class='bx bx-bar-chart-alt-2'></i>", label: 'Dashboard' },
        { href: basePath + 'trips/index.html', icon: "<i class='bx bx-mountain'></i>", label: 'Manajemen Trip' },
        { href: basePath + 'bookings/index.html', icon: "<i class='bx bx-file'></i>", label: 'Manajemen Booking' },
        { href: basePath + 'reports/revenue.html', icon: "<i class='bx bx-wallet'></i>", label: 'Laporan' },
      ]
    },
    {
      label: 'Master Data',
      items: [
        { href: basePath + 'packages/index.html', icon: "<i class='bx bx-package'></i>", label: 'Paket' },
        { href: basePath + 'meeting-points/index.html', icon: "<i class='bx bx-map-pin'></i>", label: 'Meeting Point' },
      ]
    },
    {
      label: 'Sistem',
      items: [
        { href: basePath + 'profile/index.html', icon: "<i class='bx bx-cog'></i>", label: 'Pengaturan' },
      ]
    }
  ];
  
  sidebar.innerHTML = menuSections.map(section => {
    const itemsHTML = section.items.map(item => {
      const isActive = currentPage === item.href.split('/').pop();
      return `
        <a href="${item.href}" class="sidebar-item ${isActive ? 'active' : ''}">
          <span class="icon">${item.icon}</span>
          <span>${item.label}</span>
        </a>
      `;
    }).join('');
    
    return `
      <div class="sidebar-section">${section.label}</div>
      ${itemsHTML}
    `;
  }).join('');
}

// Setup sidebar search
function setupSidebarSearch() {
  const sidebar = document.querySelector('.admin-sidebar');
  if (!sidebar) return;
  
  // Check if search already exists
  if (sidebar.querySelector('.sidebar-search')) return;
  
  const searchHTML = `
    <div class="sidebar-search">
      <input type="text" placeholder="Cari menu..." id="sidebar-search-input">
    </div>
  `;
  
  const sidebarMenu = sidebar.querySelector('.sidebar-menu');
  if (sidebarMenu) {
    sidebarMenu.insertAdjacentHTML('beforebegin', searchHTML);
  }
  
  const searchInput = document.getElementById('sidebar-search-input');
  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      const query = e.target.value.toLowerCase();
      const items = sidebar.querySelectorAll('.sidebar-item');
      const sections = sidebar.querySelectorAll('.sidebar-section');
      
      items.forEach(item => {
        const text = item.textContent.toLowerCase();
        if (text.includes(query)) {
          item.style.display = '';
        } else {
          item.style.display = 'none';
        }
      });
      
      // Hide/show sections based on visible items
      sections.forEach(section => {
        const nextItems = [];
        let nextEl = section.nextElementSibling;
        while (nextEl && !nextEl.classList.contains('sidebar-section')) {
          if (nextEl.classList.contains('sidebar-item')) {
            nextItems.push(nextEl);
          }
          nextEl = nextEl.nextElementSibling;
        }
        
        const hasVisible = nextItems.some(item => item.style.display !== 'none');
        section.style.display = hasVisible ? '' : 'none';
      });
    });
  }
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

// Get operator meeting points
function getOperatorMeetingPoints() {
  return getMeetingPointsByOperator(CURRENT_OPERATOR_ID);
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

// Format date
function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' });
}

// Get status badge HTML
function getStatusBadge(status) {
  const badges = {
    pending: '<span class="badge badge-pending"><i class="bx bx-time"></i> Pending</span>',
    confirmed: '<span class="badge badge-confirmed"><i class="bx bx-check-circle"></i> Confirmed</span>',
    completed: '<span class="badge badge-completed"><i class="bx bx-check"></i> Completed</span>',
    cancelled: '<span class="badge badge-cancelled"><i class="bx bx-x"></i> Cancelled</span>'
  };
  return badges[status] || status;
}

// Logout
function logout() {
  const basePath = getBasePath();
  window.location.href = basePath + 'login.html';
}

// Initialize when DOM ready
document.addEventListener('DOMContentLoaded', initAdmin);
