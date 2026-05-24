// Admin shell utilities
document.addEventListener('DOMContentLoaded', function() {
    // Close sidebar on click outside (mobile)
    document.addEventListener('click', function(e) {
        var sidebar = document.getElementById('sidebar');
        if (sidebar && sidebar.classList.contains('open')) {
            var isSidebar = sidebar.contains(e.target);
            var isToggle = e.target.closest('[onclick="toggleSidebar()"]');
            if (!isSidebar && !isToggle) {
                sidebar.classList.remove('open');
            }
        }
    });
});
