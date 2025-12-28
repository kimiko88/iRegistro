import { mount } from '@vue/test-utils';
import { describe, it, expect } from 'vitest';
import DataTable from '@/components/shared/DataTable.vue';

describe('DataTable.vue', () => {
    const columns = [
        { key: 'name', label: 'Name', sortable: true },
        { key: 'role', label: 'Role', sortable: false },
    ];

    const data = [
        { id: 1, name: 'Alice', role: 'Admin' },
        { id: 2, name: 'Bob', role: 'User' },
    ];

    it('renders table headers correctly', () => {
        const wrapper = mount(DataTable, {
            props: { columns, data }
        });
        const headers = wrapper.findAll('th');
        expect(headers).toHaveLength(columns.length + 1); // +1 for Actions column by default if hasActions is true (which it is hardcoded to be currently)
        expect(headers[0].text()).toContain('Name');
        expect(headers[1].text()).toContain('Role');
    });

    it('renders data rows correctly', () => {
        const wrapper = mount(DataTable, {
            props: { columns, data }
        });
        const rows = wrapper.findAll('tbody tr');
        // data has 2 items, so 2 rows.
        expect(rows).toHaveLength(2);
        expect(rows[0].text()).toContain('Alice');
        expect(rows[1].text()).toContain('Bob');
    });

    it('emits search event when input changes', async () => {
        const wrapper = mount(DataTable, {
            props: { columns, data, enableSearch: true }
        });
        const input = wrapper.find('input[type="text"]');
        await input.setValue('test');
        expect(wrapper.emitted()).toHaveProperty('search');
        expect(wrapper.emitted('search')![0]).toEqual(['test']);
    });

    it('emits sort event when header is clicked', async () => {
        const wrapper = mount(DataTable, {
            props: { columns, data }
        });
        const nameHeader = wrapper.find('th'); // first one is Name, sortable
        await nameHeader.trigger('click');
        expect(wrapper.emitted()).toHaveProperty('update:sort');
        expect(wrapper.emitted('update:sort')![0]).toEqual([{ key: 'name', order: 'asc' }]);

        await nameHeader.trigger('click');
        expect(wrapper.emitted('update:sort')![1]).toEqual([{ key: 'name', order: 'desc' }]);
    });
});
