<script setup lang="ts">
// Student dashboard logic (uses same generic views but different data source potentially)
</script>

<template>
  <div class="p-4">
      <h1 class="text-2xl font-bold mb-6">My Dashboard</h1>
      <div class="hero bg-base-200 rounded-lg p-6">
          <div class="hero-content text-center">
              <div class="max-w-md">
                  <h1 class="text-5xl font-bold">Hello Student</h1>
                  <p class="py-6">Here is your overview.</p>
                  
                  <div class="stats shadow">
                    <div class="stat">
                        <div class="stat-title">Average</div>
                        <div class="stat-value">8.5</div>
                    </div>
                    <div class="stat">
                        <div class="stat-title">Absences</div>
                        <div class="stat-value">3</div>
                    </div>
                     <div class="stat">
                        <div class="stat-title">Next</div>
                        <div class="stat-value text-sm">Math Test</div>
                    </div>
                  </div>

                  <button class="btn btn-primary mt-4" @click="$router.push('/student/marks')">View Grades</button>
              </div>
          </div>
      </div>
      <!-- Add generic components here -->
  </div>
</template>
