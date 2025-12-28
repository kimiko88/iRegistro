<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
    marks: any[];
    average?: number;
}>();

const groupBySubject = computed(() => {
    const groups: Record<string, any[]> = {};
    props.marks.forEach(m => {
        if (!groups[m.subject]) groups[m.subject] = [];
        groups[m.subject].push(m);
    });
    return groups;
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center bg-base-100 p-4 rounded shadow">
        <h2 class="text-xl font-bold">Grades Analysis</h2>
        <div class="stats stats-horizontal shadow">
            <div class="stat">
                <div class="stat-title">GPA</div>
                <div class="stat-value text-primary">{{ average || '-' }}</div>
            </div>
        </div>
    </div>

    <!-- Timeline / List -->
    <div class="card bg-base-100 shadow">
        <div class="card-body">
            <h3 class="card-title">Recent Marks</h3>
            <div class="overflow-x-auto">
                <table class="table table-zebra w-full">
                    <thead>
                        <tr>
                            <th>Date</th>
                            <th>Subject</th>
                            <th>Value</th>
                            <th>Type</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="mark in marks" :key="mark.id">
                            <td>{{ mark.date }}</td>
                            <td class="font-bold">{{ mark.subject }}</td>
                            <td>
                                <div class="badge badge-lg" :class="mark.value >= 6 ? 'badge-success' : 'badge-error'">
                                    {{ mark.value }}
                                </div>
                            </td>
                            <td>{{ mark.type || 'Oral' }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
  </div>
</template>
