import api from './api';

export default {
    // Notifications
    getNotifications() {
        return api.get('/communication/notifications');
    },
    markNotificationRead(id: number) {
        return api.post(`/communication/notifications/${id}/read`);
    },

    // Messaging
    getConversations() {
        return api.get('/communication/conversations');
    },
    getMessages(conversationId: number) {
        return api.get(`/communication/conversations/${conversationId}/messages`);
    },
    sendMessage(conversationId: number, body: string, attachments: any[] = []) {
        return api.post(`/communication/conversations/${conversationId}/messages`, { body, attachments });
    },
    createConversation(recipientId: number, subject: string, body: string) {
        return api.post('/communication/conversations', { recipient_id: recipientId, subject, body });
    },

    // Colloquiums
    getAvailableSlots(teacherId?: number) {
        return api.get('/communication/slots/available', { params: { teacher_id: teacherId } });
    },
    bookSlot(slotId: number, notes: string = '') {
        return api.post('/communication/bookings', { slot_id: slotId, notes_before: notes });
    }
};
