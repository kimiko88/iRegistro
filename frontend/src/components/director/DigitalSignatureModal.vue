<template>
    <dialog class="modal" :class="{ 'modal-open': isOpen }">
        <div class="modal-box">
            <h3 class="font-bold text-lg">Digital Signature Required</h3>
            <p class="py-4">Please enter your 6-digit PIN to sign this document digitally.</p>
            
            <div class="form-control w-full max-w-xs mx-auto">
                <label class="label">
                    <span class="label-text">Security PIN</span>
                </label>
                <input 
                    type="password" 
                    placeholder="123456" 
                    class="input input-bordered w-full max-w-xs text-center tracking-widest text-xl" 
                    v-model="pin"
                    maxlength="6"
                />
                <label class="label" v-if="error">
                    <span class="label-text-alt text-error">{{ error }}</span>
                </label>
            </div>

            <div class="modal-action">
                <button class="btn" @click="cancel" :disabled="loading">Cancel</button>
                <button class="btn btn-primary" @click="confirm" :disabled="pin.length !== 6 || loading">
                    <span v-if="loading" class="loading loading-spinner"></span>
                    Sign Document
                </button>
            </div>
        </div>
    </dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
    isOpen: boolean;
    loading?: boolean;
}>();

const emit = defineEmits(['confirm', 'cancel']);

const pin = ref('');
const error = ref('');

watch(() => props.isOpen, (newVal) => {
    if (newVal) {
        pin.value = '';
        error.value = '';
    }
});

const cancel = () => {
    emit('cancel');
};

const confirm = () => {
    if (pin.value.length !== 6) {
        error.value = 'PIN must be 6 digits';
        return;
    }
    emit('confirm', pin.value);
};
</script>
