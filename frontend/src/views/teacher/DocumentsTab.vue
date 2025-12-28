<template>
  <div class="h-full flex flex-col p-4 bg-base-100 rounded-box shadow gap-6 overflow-y-auto">
      <!-- Section: General Documents -->
      <div>
        <h2 class="text-xl font-bold mb-4">Class Documents</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <!-- Document Cards -->
            <div class="card bg-base-100 shadow border border-base-200 hover:border-primary transition-colors cursor-pointer group">
                <div class="card-body items-center text-center p-6">
                    <div class="text-4xl mb-2 group-hover:scale-110 transition-transform">📄</div>
                    <h3 class="card-title text-sm">Documento 15 Maggio</h3>
                    <p class="text-xs opacity-60">Scadenza: 15/05/2024</p>
                    <div class="card-actions mt-4">
                        <button class="btn btn-sm btn-primary">Compila</button>
                    </div>
                </div>
            </div>
            
             <div class="card bg-base-100 shadow border border-base-200 hover:border-primary transition-colors cursor-pointer group">
                <div class="card-body items-center text-center p-6">
                    <div class="text-4xl mb-2 group-hover:scale-110 transition-transform">📊</div>
                    <h3 class="card-title text-sm">Programma Svolto</h3>
                    <p class="text-xs opacity-60">Ultimo agg: Ieri</p>
                     <div class="card-actions mt-4">
                        <button class="btn btn-sm btn-outline">Vedi</button>
                    </div>
                </div>
            </div>

             <div class="card bg-base-100 shadow border border-base-200 hover:border-primary transition-colors cursor-pointer group">
                <div class="card-body items-center text-center p-6">
                    <div class="text-4xl mb-2 group-hover:scale-110 transition-transform">🎓</div>
                    <h3 class="card-title text-sm">Pagelle (Draft)</h3>
                    <p class="text-xs opacity-60">1Q Term</p>
                     <div class="card-actions mt-4">
                        <button class="btn btn-sm btn-outline">Review</button>
                    </div>
                </div>
            </div>
        </div>
      </div>

      <div class="divider"></div>

      <!-- Section: Student Specific -->
      <div>
        <h2 class="text-xl font-bold mb-4">Student Documents (PDP/PEI)</h2>
        <div class="overflow-x-auto border border-base-200 rounded-box">
            <table class="table w-full">
                <thead>
                    <tr class="bg-base-200/50">
                        <th>Student</th>
                        <th>Document Type</th>
                        <th>Created At</th>
                        <th>Status</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody>
                    <tr class="hover">
                        <td class="font-bold">Rossi Luigi</td>
                        <td>PDP (DSA)</td>
                        <td class="text-sm opacity-70">10/10/2023</td>
                        <td><span class="badge badge-success gap-1">Signed</span></td>
                        <td><button class="btn btn-xs btn-ghost text-primary">Download</button></td>
                    </tr>
                     <tr class="hover">
                        <td class="font-bold">Bianchi Marco</td>
                        <td>PEI</td>
                        <td class="text-sm opacity-70">12/11/2023</td>
                        <td><span class="badge badge-warning gap-1">Pending Sign</span></td>
                        <td>
                            <button class="btn btn-xs btn-primary gap-1">
                                Sign
                            </button>
                        </td>
                    </tr>
                    <!-- Mock more rows -->
                    <tr v-if="false"></tr>
                </tbody>
            </table>
        </div>
      </div>
  </div>
</template>

<script setup lang="ts">
// Logic for fetching documents goes here
</script>
