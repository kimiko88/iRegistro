import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import RequestApprovals from '../RequestApprovals.vue';

const approveSpy = vi.fn();
const rejectSpy = vi.fn();
const fetchRequestsSpy = vi.fn();

vi.mock('@/stores/director', () => ({
    useDirectorStore: () => ({
        pendingRequests: [
            { id: 101, type: 'TRANSFER', requester: 'Parent A', details: 'Moving', date: '2025-02-02' }
        ],
        approveRequest: approveSpy,
        rejectRequest: rejectSpy,
        fetchRequests: fetchRequestsSpy
    })
}));

describe('RequestApprovals.vue', () => {
    it('renders requests and actions', async () => {
        const wrapper = mount(RequestApprovals);

        expect(fetchRequestsSpy).toHaveBeenCalled();
        expect(wrapper.text()).toContain('TRANSFER');
        expect(wrapper.text()).toContain('Parent A');

        const approveBtn = wrapper.findAll('button').find(b => b.text() === 'Approve');
        await approveBtn?.trigger('click');
        expect(approveSpy).toHaveBeenCalledWith(101);

        const rejectBtn = wrapper.findAll('button').find(b => b.text() === 'Reject');
        await rejectBtn?.trigger('click');
        expect(rejectSpy).toHaveBeenCalledWith(101);
    });
});
