<template>
  <div class="admin-container">
    <h2>Admin Control Panel</h2>
    
    <div class="admin-grid">
      <router-link to="/list-events" class="admin-card">
        <div class="admin-card-icon">📅</div>
        <h3>Events</h3>
        <p>View and manage events</p>
      </router-link>

      <router-link to="/diagnostic" class="admin-card">
        <div class="admin-card-icon">🔧</div>
        <h3>Diagnostics</h3>
        <p>CBOR decoder and diagnostic tools</p>
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Admin'
}
</script>

<style scoped>
.admin-container {
  max-width: 1200px;
  margin: 0 auto;
}

h2 {
  margin-bottom: 2rem;
  color: #2c3e50;
  font-size: 2rem;
}

.admin-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.admin-card {
  background: white;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  padding: 2rem;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.admin-card:hover {
  border-color: #3498db;
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.15);
  transform: translateY(-2px);
}

.admin-card-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.admin-card h3 {
  margin: 0.5rem 0;
  color: #2c3e50;
  font-size: 1.5rem;
}

.admin-card p {
  margin: 0;
  color: #666;
  font-size: 0.95rem;
}
</style>

