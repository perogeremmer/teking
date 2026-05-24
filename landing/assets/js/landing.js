document.addEventListener('DOMContentLoaded', function () {
  var navbar = document.getElementById('navbar');
  var navInner = document.getElementById('nav-inner');
  var menuToggle = document.getElementById('menu-toggle');
  var menuClose = document.getElementById('menu-close');
  var mobileMenu = document.getElementById('mobile-menu');
  var menuOverlay = document.getElementById('menu-overlay');

  function isDesktop() {
    return window.innerWidth >= 1024;
  }

  // Navbar shrink on scroll
  function handleScroll() {
    var scrolled = window.scrollY > 50;
    navbar.classList.toggle('shadow-sm', scrolled);
    if (navInner) {
      navInner.classList.toggle('h-16', !scrolled);
      navInner.classList.toggle('lg:h-18', !scrolled);
      navInner.classList.toggle('h-14', scrolled);
      navInner.classList.toggle('lg:h-16', scrolled);
    }
  }

  handleScroll();
  window.addEventListener('scroll', handleScroll, { passive: true });

  // Mobile menu
  function openMenu() {
    mobileMenu.classList.remove('translate-x-full');
    menuOverlay.classList.remove('hidden');
    document.body.style.overflow = 'hidden';
  }

  function closeMenu() {
    mobileMenu.classList.add('translate-x-full');
    menuOverlay.classList.add('hidden');
    document.body.style.overflow = '';
  }

  if (menuToggle) menuToggle.addEventListener('click', openMenu);
  if (menuClose) menuClose.addEventListener('click', closeMenu);
  if (menuOverlay) menuOverlay.addEventListener('click', closeMenu);

  mobileMenu.querySelectorAll('a').forEach(function (link) {
    link.addEventListener('click', closeMenu);
  });

  window.addEventListener('resize', function () {
    if (isDesktop()) closeMenu();
  });

  // IntersectionObserver for reveal animations
  var observer = new IntersectionObserver(function (entries) {
    entries.forEach(function (entry) {
      if (entry.isIntersecting) {
        var target = entry.target;

        if (target.classList.contains('reveal')) {
          target.classList.add('visible');
          observer.unobserve(target);
          return;
        }

        var children = target.querySelectorAll('.reveal-child');
        if (children.length > 0) {
          observer.unobserve(target);
          target.querySelectorAll('.reveal').forEach(function (el) {
            el.classList.add('visible');
          });
          children.forEach(function (child, i) {
            setTimeout(function () {
              child.classList.add('visible');
            }, i * 120);
          });
          return;
        }
      }
    });
  }, {
    threshold: 0.08,
    rootMargin: '0px 0px -40px 0px'
  });

  document.querySelectorAll('.reveal, [class*="reveal-child"]').forEach(function (el) {
    var parent = el.closest('section, footer, header');
    if (parent && (el.classList.contains('reveal-child') || parent.querySelector('.reveal-child'))) {
      if (!parent.dataset.observed) {
        parent.dataset.observed = 'true';
        observer.observe(parent);
      }
    } else if (el.classList.contains('reveal')) {
      observer.observe(el);
    }
  });

  // Smooth scroll for anchor links
  document.querySelectorAll('a[href^="#"]').forEach(function (link) {
    link.addEventListener('click', function (e) {
      var targetId = this.getAttribute('href');
      if (targetId === '#') return;
      var target = document.querySelector(targetId);
      if (target) {
        e.preventDefault();
        var navHeight = navbar ? navbar.offsetHeight : 0;
        var targetPos = target.getBoundingClientRect().top + window.pageYOffset - navHeight;
        window.scrollTo({ top: targetPos, behavior: 'smooth' });
      }
    });
  });
});