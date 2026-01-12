<template>
  <div class="app-container">
    <header class="app-header">
      <h1>AttendeeAttestor</h1>
      <button @click="toggleFullscreen" class="fullscreen-button" :title="isFullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'">
        <span v-if="!isFullscreen">⛶</span>
        <span v-else>⛶</span>
      </button>
    </header>
    <nav class="navigation">
      <router-link
        to="/scanner"
        class="nav-link"
        active-class="active"
      >
        Scanner
      </router-link>
      <router-link
        to="/admin"
        class="nav-link"
        :class="{ active: isAdminRoute }"
      >
        Admin
      </router-link>
    </nav>
    <main class="app-main">
      <div class="content-sections">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      isFullscreen: false
    }
  },
  computed: {
    isAdminRoute() {
      const adminRoutes = ['/admin', '/list-events', '/create-event', '/issue-ticket', '/issued-tickets', '/diagnostic', '/ticket-details']
      return adminRoutes.includes(this.$route.path)
    }
  },
  mounted() {
    document.addEventListener('fullscreenchange', this.handleFullscreenChange)
    document.addEventListener('webkitfullscreenchange', this.handleFullscreenChange)
    document.addEventListener('mozfullscreenchange', this.handleFullscreenChange)
    document.addEventListener('MSFullscreenChange', this.handleFullscreenChange)
  },
  beforeUnmount() {
    document.removeEventListener('fullscreenchange', this.handleFullscreenChange)
    document.removeEventListener('webkitfullscreenchange', this.handleFullscreenChange)
    document.removeEventListener('mozfullscreenchange', this.handleFullscreenChange)
    document.removeEventListener('MSFullscreenChange', this.handleFullscreenChange)
  },
  methods: {
    toggleFullscreen() {
      if (!document.fullscreenElement && !document.webkitFullscreenElement && !document.mozFullScreenElement && !document.msFullscreenElement) {
        const element = document.documentElement
        if (element.requestFullscreen) {
          element.requestFullscreen()
        } else if (element.webkitRequestFullscreen) {
          element.webkitRequestFullscreen()
        } else if (element.mozRequestFullScreen) {
          element.mozRequestFullScreen()
        } else if (element.msRequestFullscreen) {
          element.msRequestFullscreen()
        }
      } else {
        if (document.exitFullscreen) {
          document.exitFullscreen()
        } else if (document.webkitExitFullscreen) {
          document.webkitExitFullscreen()
        } else if (document.mozCancelFullScreen) {
          document.mozCancelFullScreen()
        } else if (document.msExitFullscreen) {
          document.msExitFullscreen()
        }
      }
    },
    handleFullscreenChange() {
      this.isFullscreen = !!(
        document.fullscreenElement ||
        document.webkitFullscreenElement ||
        document.mozFullScreenElement ||
        document.msFullscreenElement
      )
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background-color: #f5f5f5;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background-color: #2c3e50;
  color: white;
  padding: 1.5rem 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-header h1 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.fullscreen-button {
  background-color: transparent;
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  padding: 0.5rem 0.75rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1.25rem;
  line-height: 1;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 2.5rem;
  height: 2.5rem;
}

.fullscreen-button:hover {
  background-color: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.5);
}

.fullscreen-button:active {
  background-color: rgba(255, 255, 255, 0.2);
}

.navigation {
  background-color: white;
  border-bottom: 1px solid #ddd;
  padding: 0 2rem;
  display: flex;
  gap: 0.5rem;
}

.nav-link {
  padding: 1rem 1.5rem;
  background-color: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  color: #666;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  display: inline-block;
  transition: all 0.2s;
}

.nav-link:hover {
  color: #2c3e50;
  background-color: #f5f5f5;
}

.nav-link.active {
  color: #3498db;
  border-bottom-color: #3498db;
}

.app-main {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
}

.content-sections {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
</style>

