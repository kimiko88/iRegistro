<template>
    <div class="p-6">
        <h1 class="text-3xl font-bold mb-6">School Reports & Analytics</h1>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <!-- Chart Section (Mock Visuals) -->
            <div class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <h2 class="card-title">Average Grade Distribution</h2>
                    <div class="h-64 flex items-end justify-between px-4 pt-4 border-b border-l border-gray-300 gap-2">
                        <!-- Mock bars -->
                         <div class="w-full bg-blue-100 hover:bg-blue-200 rounded-t relative group h-10"></div>
                         <div class="w-full bg-blue-200 hover:bg-blue-300 rounded-t relative group h-24"></div>
                         <div class="w-full bg-blue-300 hover:bg-blue-400 rounded-t relative group h-40"></div>
                         <div class="w-full bg-blue-500 hover:bg-blue-600 rounded-t relative group h-56"></div>
                         <div class="w-full bg-blue-400 hover:bg-blue-500 rounded-t relative group h-32"></div>
                    </div>
                    <div class="flex justify-between text-xs text-gray-500 mt-2 px-1">
                        <span>&lt;5</span><span>5-6</span><span>6-7</span><span>7-8</span><span>8-10</span>
                    </div>
                </div>
            </div>

             <div class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <h2 class="card-title">Absence Anomalies</h2>
                    <table class="table table-compact w-full">
                        <thead>
                            <tr>
                                <th>Class</th>
                                <th>Anomaly Rate</th>
                                <th>Status</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td>3C</td>
                                <td class="text-red-500 font-bold">+15% vs Avg</td>
                                <td><div class="badge badge-error">Critical</div></td>
                            </tr>
                            <tr>
                                <td>5B</td>
                                <td class="text-orange-500 font-bold">+8% vs Avg</td>
                                <td><div class="badge badge-warning">Warning</div></td>
                            </tr>
                            <tr>
                                <td>1A</td>
                                <td class="text-green-500">-2% vs Avg</td>
                                <td><div class="badge badge-success">OK</div></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            
            <!-- Generator Component -->
            <div class="col-span-1 lg:col-span-2">
                <ReportGenerator />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
// @ts-ignore
import ReportGenerator from '@/components/director/ReportGenerator.vue';
</script>
