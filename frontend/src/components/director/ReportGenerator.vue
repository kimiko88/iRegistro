<template>
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h2 class="card-title">Generate Report</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Report Type</span>
                    </label>
                    <select class="select select-bordered" v-model="reportType">
                        <option value="attendance">Class Attendance</option>
                        <option value="grades">Grade Distribution</option>
                        <option value="teacher_perf">Teacher Performance</option>
                        <option value="year_comparison">Yearly Comparison</option>
                    </select>
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Format</span>
                    </label>
                    <select class="select select-bordered" v-model="format">
                        <option value="pdf">PDF Document</option>
                        <option value="csv">CSV Export</option>
                        <option value="excel">Excel Spreadsheet</option>
                    </select>
                </div>
                
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Academic Year</span>
                    </label>
                    <select class="select select-bordered" v-model="year">
                        <option>2024-2025</option>
                        <option>2023-2024</option>
                    </select>
                </div>

                <div class="form-control mt-9">
                    <button class="btn btn-secondary w-full" @click="generate">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 mr-2">
                           <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
                         </svg>
                        Download Report
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const reportType = ref('attendance');
const format = ref('pdf');
const year = ref('2024-2025');

const generate = () => {
    alert(`Generating ${reportType.value} report in ${format.value} format for ${year.value}... (Mock)`);
};
</script>
