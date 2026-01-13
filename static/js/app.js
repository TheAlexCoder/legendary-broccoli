// DOM ready function to initialize Bootstrap modal
document.addEventListener('DOMContentLoaded', function() {
    // Initialize the modal element
    const recoveryModalElement = document.getElementById('recoveryPhraseModal');
    if (recoveryModalElement) {
        window.recoveryPhraseModal = new bootstrap.Modal(recoveryModalElement);
        
        // Add event listener for copy button
        document.getElementById('copyRecoveryBtn').addEventListener('click', function() {
            const phrase = document.getElementById('recoveryPhraseDisplay').textContent;
            navigator.clipboard.writeText(phrase).then(() => {
                alert('Recovery phrase copied to clipboard!');
            }).catch(err => {
                console.error('Failed to copy: ', err);
            });
        });
        
        // Add event listener for when modal is hidden
        recoveryModalElement.addEventListener('hidden.bs.modal', function () {
            localStorage.setItem('hasClosedRecoveryModal', 'true');
        });
    }
});

function appData() {
    return {
        isLoggedIn: false,
        showLogin: true,
        regUsername: '',
        loginRecoveryPhrase: '',
        recoveryPhrase: '',
        username: '',
        firedDate: '',
        stats: {
            total_working_days: 0,
            checked_days: 0,
            remaining_days: 0
        },
        calendarDays: [],
        currentMonth: new Date().getMonth(),
        currentYear: new Date().getFullYear(),
        showSettings: false,
        showRecoveryPhrase: false,
        currentLang: 'en',
        translations: {
            en: enTranslations,
            ru: ruTranslations,
        },
        
        init() {
            // Initialize language from localStorage or default to 'en'
            const savedLang = localStorage.getItem('lang');
            if (savedLang && (savedLang === 'en' || savedLang === 'ru')) {
                this.currentLang = savedLang;
            } else {
                // Try to detect browser language
                const browserLang = navigator.language.split('-')[0];
                if (browserLang === 'ru') {
                    this.currentLang = 'ru';
                } else {
                    this.currentLang = 'en';
                }
                localStorage.setItem('lang', this.currentLang);
            }
            
            // Check if user is already logged in by making a request to get profile
            this.checkAuthStatus();
        },
        
        async checkAuthStatus() {
            try {
                const response = await fetch('/api/profile');
                if (response.ok) {
                    const data = await response.json();
                    this.isLoggedIn = true;
                    this.username = data.username;
                    this.firedDate = data.fired_date || '';
                    
                    // Check if fired date is not set and show notification
                    if (!this.firedDate) {
                        this.showNotification(this.t('firedDateNotSet') || 'Please set your fired date in settings.', 'warning');
                    }
                    
                    // Load calendar and stats
                    await this.loadCalendar();
                    await this.loadStats();
                } else {
                    this.isLoggedIn = false;
                }
            } catch (error) {
                console.error('Error checking auth status:', error);
                this.isLoggedIn = false;
            }
        },
        
        async register() {
            try {
                const response = await fetch('/api/auth/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        username: this.regUsername
                    })
                });
                
                if (response.ok) {
                    const data = await response.json();
                    this.recoveryPhrase = data.recovery_phrase;
                    
                    // Auto-login by using the recovery phrase
                    await this.loginWithRecoveryPhrase(data.recovery_phrase);
                    
                    // Show modal with recovery phrase if user hasn't closed it before
                    const hasClosedModal = localStorage.getItem('hasClosedRecoveryModal');
                    if (!hasClosedModal) {
                        document.getElementById('recoveryPhraseDisplay').textContent = this.recoveryPhrase;
                        if (window.recoveryPhraseModal) {
                            window.recoveryPhraseModal.show();
                        }
                    }
                    
                    // Reset form
                    this.regUsername = '';
                } else {
                    const error = await response.json();
                    alert('Registration failed: ' + (error.message || 'Unknown error'));
                }
            } catch (error) {
                console.error('Registration error:', error);
                alert('Registration failed: ' + error.message);
            }
        },
        
        async login() {
            try {
                await this.loginWithRecoveryPhrase(this.loginRecoveryPhrase);
            } catch (error) {
                console.error('Login error:', error);
                alert('Login failed: ' + error.message);
            }
        },
        
        async loginWithRecoveryPhrase(recoveryPhrase) {
            try {
                const response = await fetch('/api/auth/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        recovery_phrase: recoveryPhrase
                    })
                });
                
                if (response.ok) {
                    this.isLoggedIn = true;
                    this.loginRecoveryPhrase = '';
                    
                    // Load user data
                    await this.checkAuthStatus();
                } else {
                    const error = await response.json();
                    throw new Error(error.message || 'Invalid recovery phrase');
                }
            } catch (error) {
                throw error; // Re-throw to be handled by caller
            }
        },
        
        async logout() {
            try {
                await fetch('/api/auth/logout', {
                    method: 'POST'
                });
                
                this.isLoggedIn = false;
                this.username = '';
                this.recoveryPhrase = '';
                this.stats = { total_working_days: 0, checked_days: 0, remaining_days: 0 };
                this.calendarDays = [];
            } catch (error) {
                console.error('Logout error:', error);
            }
        },
        
        async loadStats() {
            try {
                const response = await fetch('/api/calendar/stats');
                if (response.ok) {
                    this.stats = await response.json();
                } else if (response.status === 400) {
                    // Handle case where fired date is not set
                    const error = await response.json();
                    if (error.message && error.message.includes('Fired date not set')) {
                        // Show notification that fired date needs to be set
                        this.showNotification('Please set your fired date in settings.', 'warning');
                    }
                }
            } catch (error) {
                console.error('Error loading stats:', error);
            }
        },
        
        async loadCalendar() {
            try {
                const year = this.currentYear;
                const month = String(this.currentMonth + 1).padStart(2, '0');
                
                // Get first and last day of the month
                const firstDay = new Date(year, this.currentMonth, 1);
                const lastDay = new Date(year, this.currentMonth + 1, 0);
                
                const response = await fetch(`/api/calendar/days`);
                if (response.ok) {
                    const entries = await response.json();
                    const entriesMap = {};
                    if (entries && Array.isArray(entries)) {
                        entries.forEach(entry => {
                            const date = new Date(entry.date).toISOString().split('T')[0]
                            entriesMap[date] = entry.checked;
                        });
                    }
                    
                    this.generateCalendarDays(firstDay, lastDay, entriesMap);
                }
            } catch (error) {
                console.error('Error loading calendar:', error);
            }
        },
        
        generateCalendarDays(firstDay, lastDay, entriesMap) {
            const daysInMonth = lastDay.getDate();
            const startDay = firstDay.getDay(); // 0 = Sunday, 1 = Monday, etc.
            
            this.calendarDays = [];
            
            // Create a matrix of 6 weeks x 7 days (max 42 days)
            let dayCount = 0;
            
            // Calculate the first date to show (may be from previous month)
            const startDate = new Date(firstDay);
            startDate.setDate(startDate.getDate() - startDay);
            
            // Create calendar grid - 6 weeks maximum
            for (let weekIndex = 0; weekIndex < 6; weekIndex++) {
                const week = [];
                
                for (let dayIndex = 0; dayIndex < 7; dayIndex++) {
                    const currentDate = new Date(startDate);
                    currentDate.setDate(startDate.getDate() + (weekIndex * 7) + dayIndex);
                    
                    const dateString = currentDate.toISOString().split('T')[0];
                    const isCurrentMonth = currentDate.getMonth() === this.currentMonth &&
                                          currentDate.getFullYear() === this.currentYear;
                    
                    const dayObj = {
                        date: currentDate,
                        isCurrentMonth: isCurrentMonth,
                        checked: entriesMap[dateString] || false,
                        isPast: currentDate < new Date(),
                        isFuture: currentDate > new Date()
                    };
                    
                    week.push(dayObj);
                }
                
                this.calendarDays.push(week);
            }
        },
        
        async toggleDay(day) {
            if (!day.isCurrentMonth) return;
            
            try {
                if (day.checked) {
                    // Uncheck the day
                    const response = await fetch('/api/calendar/uncheck', {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            date: day.date.toISOString().split('T')[0]
                        })
                    });
                    
                    if (response.ok) {
                        day.checked = false;
                        await this.loadStats(); // Refresh stats after toggling
                    }
                } else {
                    // Check the day
                    const response = await fetch('/api/calendar/check', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            date: day.date.toISOString().split('T')[0]
                        })
                    });
                    
                    if (response.ok) {
                        day.checked = true;
                        await this.loadStats(); // Refresh stats after toggling
                    }
                }
            } catch (error) {
                console.error('Error toggling day:', error);
            }
        },
        
        async markToday() {
            const today = new Date();
            const todayStr = today.toISOString().split('T')[0];
            
            // Find today in the calendar
            for (const week of this.calendarDays) {
                for (const day of week) {
                    if (day && day.date.toISOString().split('T')[0] === todayStr) {
                        // Toggle today's status
                        await this.toggleDay(day);
                        return;
                    }
                }
            }
        },
        
        async updateProfile() {
            try {
                const response = await fetch('/api/profile', {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        username: this.username,
                        fired_date: this.firedDate
                    })
                });
                
                if (!response.ok) {
                    const error = await response.json();
                    alert('Update failed: ' + (error.message || 'Unknown error'));
                }
            } catch (error) {
                console.error('Update error:', error);
            }
        },
        
        async showRecoveryPhraseModal() {
            try {
                const response = await fetch('/api/profile/recovery');
                if (response.ok) {
                    const data = await response.json();
                    this.recoveryPhrase = data.recovery_phrase;
                    
                    // Show modal with recovery phrase
                    document.getElementById('recoveryPhraseDisplay').textContent = this.recoveryPhrase;
                    if (window.recoveryPhraseModal) {
                        window.recoveryPhraseModal.show();
                    }
                } else {
                    alert('Failed to load recovery phrase');
                }
            } catch (error) {
                console.error('Error loading recovery phrase:', error);
                alert('Error loading recovery phrase');
            }
        },
        
        async loadRecoveryPhrase() {
            try {
                const response = await fetch('/api/profile/recovery');
                if (response.ok) {
                    const data = await response.json();
                    this.recoveryPhrase = data.recovery_phrase;
                }
            } catch (error) {
                console.error('Error loading recovery phrase:', error);
            }
        },
        
        async deleteAccount() {
            if (!confirm('Are you sure you want to delete your account? It can be restored within 7 days.')) {
                return;
            }
            
            try {
                const response = await fetch('/api/user/delete', {
                    method: 'DELETE'
                });
                
                if (response.ok) {
                    alert('Account deleted. You can restore it within 7 days.');
                    await this.logout();
                } else {
                    const error = await response.json();
                    alert('Delete failed: ' + (error.message || 'Unknown error'));
                }
            } catch (error) {
                console.error('Delete error:', error);
            }
        },
        
        copyRecoveryPhrase() {
            navigator.clipboard.writeText(this.recoveryPhrase).then(() => {
                alert('Recovery phrase copied to clipboard!');
            }).catch(err => {
                console.error('Failed to copy: ', err);
            });
        },
        
        get currentMonthYear() {
            const months = [
                'January', 'February', 'March', 'April', 'May', 'June',
                'July', 'August', 'September', 'October', 'November', 'December'
            ];
            return `${months[this.currentMonth]} ${this.currentYear}`;
        },
        
        async prevMonth() {
            if (this.currentMonth === 0) {
                this.currentMonth = 11;
                this.currentYear--;
            } else {
                this.currentMonth--;
            }
            await this.loadCalendar();
        },
        
        async nextMonth() {
            if (this.currentMonth === 11) {
                this.currentMonth = 0;
                this.currentYear++;
            } else {
                this.currentMonth++;
            }
            await this.loadCalendar();
        },
        
        t(key) {
            return this.translations[this.currentLang][key] || key;
        },
        
        switchLanguage(lang) {
            this.currentLang = lang;
            localStorage.setItem('lang', lang);
        },
        
        get currentLangName() {
            return this.currentLang === 'en' ? 'EN' : 'RU';
        },
        
        isWeekend(date) {
            const dayOfWeek = date.getDay(); // 0 = Sunday, 6 = Saturday
            return dayOfWeek === 0 || dayOfWeek === 6;
        },
        
        showNotification(message, type = 'info') {
            // Create notification element if it doesn't exist
            let notificationEl = document.getElementById('notification-toast');
            if (!notificationEl) {
                // Create toast container if it doesn't exist
                let toastContainer = document.getElementById('toast-container');
                if (!toastContainer) {
                    toastContainer = document.createElement('div');
                    toastContainer.id = 'toast-container';
                    toastContainer.className = 'toast-container position-fixed bottom-0 end-0 p-3';
                    toastContainer.style.zIndex = '1100';
                    document.body.appendChild(toastContainer);
                }
                
                notificationEl = document.createElement('div');
                notificationEl.id = 'notification-toast';
                notificationEl.className = 'toast';
                notificationEl.setAttribute('role', 'alert');
                notificationEl.innerHTML = `
                    <div class="toast-header">
                        <strong class="me-auto">Notification</strong>
                        <button type="button" class="btn-close" data-bs-dismiss="toast"></button>
                    </div>
                    <div class="toast-body" id="toast-message"></div>
                `;
                toastContainer.appendChild(notificationEl);
            }
            
            // Set message and style based on type
            const toastBody = document.getElementById('toast-message');
            toastBody.textContent = message;
            
            // Update styling based on type
            notificationEl.classList.remove('text-bg-primary', 'text-bg-success', 'text-bg-warning', 'text-bg-danger');
            switch(type) {
                case 'success':
                    notificationEl.classList.add('text-bg-success');
                    break;
                case 'warning':
                    notificationEl.classList.add('text-bg-warning');
                    break;
                case 'error':
                    notificationEl.classList.add('text-bg-danger');
                    break;
                case 'info':
                default:
                    notificationEl.classList.add('text-bg-primary');
                    break;
            }
            
            // Show the toast
            const bsToast = new bootstrap.Toast(notificationEl, { delay: 5000 });
            bsToast.show();
        }
    };
}
