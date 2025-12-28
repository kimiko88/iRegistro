import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import DocumentSigning from '../DocumentSigning.vue';
import DigitalSignatureModal from '@/components/director/DigitalSignatureModal.vue';

// Mock store
const signDocumentSpy = vi.fn();
const fetchDocumentsToSignSpy = vi.fn();

vi.mock('@/stores/director', () => ({
    useDirectorStore: () => ({
        documentsToSign: [
            { id: 1, title: 'Test Doc', type: 'Report', studentName: 'S1', date: '2025-01-01' }
        ],
        loading: false,
        signDocument: signDocumentSpy,
        fetchDocumentsToSign: fetchDocumentsToSignSpy
    })
}));

describe('DocumentSigning.vue', () => {
    it('renders document list', () => {
        const wrapper = mount(DocumentSigning);
        expect(wrapper.text()).toContain('Test Doc');
        expect(fetchDocumentsToSignSpy).toHaveBeenCalled();
    });

    it('opens modal on sign click and attempts signature', async () => {
        const wrapper = mount(DocumentSigning);

        // Find "Sign Digitally" button
        const signBtn = wrapper.find('.btn-primary');
        await signBtn.trigger('click');

        // Check if modal is visible (check class)
        const modal = wrapper.findComponent(DigitalSignatureModal);
        expect(modal.classes()).toContain('modal-open');

        // Simulate confirm event from modal
        modal.vm.$emit('confirm', '123456');

        expect(signDocumentSpy).toHaveBeenCalledWith(1, '123456');
    });
});
