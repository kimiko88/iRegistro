import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import Archive from '../Archive.vue';
import { useSecretaryStore } from '@/stores/secretary';

describe('Archive.vue', () => {
    const mocks = vi.hoisted(() => ({
        fetchArchiveSpy: vi.fn(),
    }));

    vi.mock('@/stores/secretary', () => ({
        useSecretaryStore: () => ({
            archive: [
                { id: 1, type: 'ReportCard', student: 'Mario', status: 'SIGNED', data: { file_path: 'report.pdf' } },
                { id: 2, type: 'Certificate', student: 'Luigi', status: 'DRAFT' }
            ],
            fetchArchive: mocks.fetchArchiveSpy
        })
    }));

    it('renders document list and shows download button for signed docs', async () => {
        const wrapper = mount(Archive);

        expect(mocks.fetchArchiveSpy).toHaveBeenCalled();
        const downloadBtn = wrapper.find('a[href="/api/files/download?path=report.pdf"]');
        expect(downloadBtn.exists()).toBe(true);
        expect(downloadBtn.text()).toBe('Download');

        // Check Draft doc does NOT have download button
        const rows = wrapper.findAll('tbody tr');
        expect(rows.length).toBe(2);

        const draftBtn = rows[1].find('a');
        expect(draftBtn.exists()).toBe(false);
    });
});
