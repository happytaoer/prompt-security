// API Base URL
const API_BASE = '';

// Switch between main tabs (Configuration vs Logs)
function switchTab(tabName) {
    // Show the correct main content area
    document.querySelectorAll('.tab-content').forEach(content => content.style.display = 'none');
    document.getElementById(`${tabName}-tab`).style.display = 'block';

    // Update active state for main tabs
    document.querySelectorAll('.sidebar .tab:not(.sub-tab)').forEach(tab => tab.classList.remove('active'));
    const mainTabButton = document.querySelector(`.sidebar .tab[onclick="switchTab('${tabName}')"]`);
    if (mainTabButton) mainTabButton.classList.add('active');

    // Deactivate all sub-tabs when a main tab is clicked
    document.querySelectorAll('.sub-tab').forEach(tab => tab.classList.remove('active'));

    // If switching to config, activate the first sub-tab by default
    if (tabName === 'config') {
        const firstSubTab = document.querySelector('.sidebar .sub-tab');
        if(firstSubTab) {
            firstSubTab.click();
        }
    } else if (tabName === 'logs') {
        // Deactivate all config sub-tabs when switching to logs
        document.querySelectorAll('.sidebar .sub-tab').forEach(tab => tab.classList.remove('active'));
        loadLogs();
    }
}

// Switch between configuration sections
function switchConfigSection(sectionName) {
    // Show the config tab content if it's hidden
    const configTab = document.getElementById('config-tab');
    if(configTab.style.display === 'none') {
        document.querySelectorAll('.tab-content').forEach(content => content.style.display = 'none');
        configTab.style.display = 'block';
        // Ensure the main 'Logs' tab is not active
        document.querySelector(`.sidebar .tab[onclick="switchTab('logs')"]`).classList.remove('active');
    }

    // Update active state for sub-tabs
    document.querySelectorAll('.sub-tab').forEach(tab => tab.classList.remove('active'));
    const currentSubTab = document.querySelector(`.sub-tab[onclick="switchConfigSection('${sectionName}')"]`);
    if(currentSubTab) currentSubTab.classList.add('active');

    // Show the correct config section
    document.querySelectorAll('.config-section').forEach(section => section.style.display = 'none');
    document.getElementById(`${sectionName}-section`).style.display = 'block';
}

// Load configuration from server
async function loadConfig() {
    try {
        const response = await fetch(`${API_BASE}/api/config`);
        const config = await response.json();

        // Detection settings
        document.getElementById('detect_emails').checked = config.detect_emails || false;
        document.getElementById('detect_phones').checked = config.detect_phones || false;
        document.getElementById('detect_credit_cards').checked = config.detect_credit_cards || false;
        document.getElementById('detect_ssns').checked = config.detect_ssns || false;
        document.getElementById('detect_ipv4').checked = config.detect_ipv4 || false;

        // Replacement values
        document.getElementById('email_replacement').value = config.email_replacement || '';
        document.getElementById('phone_replacement').value = config.phone_replacement || '';
        document.getElementById('credit_card_replacement').value = config.credit_card_replacement || '';
        document.getElementById('ssn_replacement').value = config.ssn_replacement || '';
        document.getElementById('ipv4_replacement').value = config.ipv4_replacement || '';

        // Monitoring settings
        document.getElementById('monitoring_interval_ms').value = config.monitoring_interval_ms || 500;
        document.getElementById('notify_on_filter').checked = config.notify_on_filter || false;

        // Custom patterns
        document.getElementById('custom_email_pattern').value = config.custom_email_pattern || '';
        document.getElementById('custom_phone_pattern').value = config.custom_phone_pattern || '';
        document.getElementById('custom_credit_card_pattern').value = config.custom_credit_card_pattern || '';
        document.getElementById('custom_ssn_pattern').value = config.custom_ssn_pattern || '';
        document.getElementById('custom_ipv4_pattern').value = config.custom_ipv4_pattern || '';

        console.log('Configuration loaded successfully');
    } catch (error) {
        console.error('Error loading configuration:', error);
        showError('Failed to load configuration');
    }
}

// Save configuration to server
async function saveConfig(event) {
    event.preventDefault();

    const config = {
        detect_emails: document.getElementById('detect_emails').checked,
        detect_phones: document.getElementById('detect_phones').checked,
        detect_credit_cards: document.getElementById('detect_credit_cards').checked,
        detect_ssns: document.getElementById('detect_ssns').checked,
        detect_ipv4: document.getElementById('detect_ipv4').checked,
        
        string_match_patterns: [], // TODO: Add UI for string patterns
        
        custom_email_pattern: document.getElementById('custom_email_pattern').value,
        custom_phone_pattern: document.getElementById('custom_phone_pattern').value,
        custom_credit_card_pattern: document.getElementById('custom_credit_card_pattern').value,
        custom_ssn_pattern: document.getElementById('custom_ssn_pattern').value,
        custom_ipv4_pattern: document.getElementById('custom_ipv4_pattern').value,
        
        email_replacement: document.getElementById('email_replacement').value,
        phone_replacement: document.getElementById('phone_replacement').value,
        credit_card_replacement: document.getElementById('credit_card_replacement').value,
        ssn_replacement: document.getElementById('ssn_replacement').value,
        ipv4_replacement: document.getElementById('ipv4_replacement').value,
        api_key_replacement: '',
        
        monitoring_interval_ms: parseInt(document.getElementById('monitoring_interval_ms').value),
        notify_on_filter: document.getElementById('notify_on_filter').checked
    };

    try {
        const response = await fetch(`${API_BASE}/api/config`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(config)
        });

        if (response.ok) {
            showSuccess('Configuration saved successfully!');
        } else {
            const error = await response.text();
            showError(`Failed to save configuration: ${error}`);
        }
    } catch (error) {
        console.error('Error saving configuration:', error);
        showError('Failed to save configuration');
    }
}

// Pagination state
let currentPage = 1;
const pageSize = 10;

// Load logs from server with pagination
async function loadLogs(page = 1) {
    try {
        const response = await fetch(`${API_BASE}/api/logs?page=${page}&pageSize=${pageSize}`);
        const data = await response.json();

        const container = document.getElementById('logs-container');
        const logs = data.logs || [];
        
        if (!logs || logs.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <p>No logs yet. Start monitoring to see filtered data.</p>
                </div>
            `;
            document.getElementById('total-logs').textContent = '0';
            document.getElementById('filtered-count').textContent = '0';
            updatePaginationButtons(1, 1);
            return;
        }

        // Update current page
        currentPage = data.page || 1;

        // Update statistics
        document.getElementById('total-logs').textContent = data.totalCount || 0;
        const totalFiltered = logs.reduce((sum, log) => sum + (log.detections?.length || 0), 0);
        document.getElementById('filtered-count').textContent = totalFiltered;

        // Render logs as table
        const tableRows = logs.map(log => {
            const timestamp = new Date(log.timestamp).toLocaleString();
            const detections = log.detections || [];
            const detectionsText = detections.length > 0 ? detections.join(', ') : '-';
            
            // Truncate text for display
            const originalText = log.original ? 
                (log.original.length > 50 ? log.original.substring(0, 50) + '...' : log.original) : 
                '-';
            const filteredText = log.filtered.length > 50 ? 
                log.filtered.substring(0, 50) + '...' : 
                log.filtered;
            
            return `
                <tr>
                    <td>${timestamp}</td>
                    <td title="${escapeHtml(log.original || '')}">${escapeHtml(originalText)}</td>
                    <td title="${escapeHtml(log.filtered)}">${escapeHtml(filteredText)}</td>
                    <td>${escapeHtml(detectionsText)}</td>
                </tr>
            `;
        }).join('');

        container.innerHTML = `
            <table class="logs-table">
                <thead>
                    <tr>
                        <th>Time</th>
                        <th>Original</th>
                        <th>Filtered</th>
                        <th>Detections</th>
                    </tr>
                </thead>
                <tbody>
                    ${tableRows}
                </tbody>
            </table>
        `;

        // Update pagination buttons
        updatePaginationButtons(currentPage, data.totalPages || 1);

    } catch (error) {
        console.error('Error loading logs:', error);
        showError('Failed to load logs');
    }
}

// Update pagination button states
function updatePaginationButtons(currentPage, totalPages) {
    const prevBtn = document.getElementById('prev-page');
    const nextBtn = document.getElementById('next-page');
    const pageInfo = document.getElementById('page-info');

    if (prevBtn) {
        prevBtn.disabled = currentPage <= 1;
    }
    if (nextBtn) {
        nextBtn.disabled = currentPage >= totalPages;
    }
    if (pageInfo) {
        pageInfo.textContent = `Page ${currentPage} / ${totalPages}`;
    }
}

// Go to previous page
function prevPage() {
    if (currentPage > 1) {
        loadLogs(currentPage - 1);
    }
}

// Go to next page
function nextPage() {
    loadLogs(currentPage + 1);
}

// Clear all logs
async function clearLogs() {
    if (!confirm('Are you sure you want to clear all logs?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/api/logs/clear`, {
            method: 'POST'
        });

        if (response.ok) {
            loadLogs();
        } else {
            showError('Failed to clear logs');
        }
    } catch (error) {
        console.error('Error clearing logs:', error);
        showError('Failed to clear logs');
    }
}

// Show success message
function showSuccess(message) {
    const element = document.getElementById('config-success');
    element.textContent = message;
    element.style.display = 'block';
    setTimeout(() => {
        element.style.display = 'none';
    }, 3000);
}

// Show error message
function showError(message) {
    const element = document.getElementById('config-error');
    element.textContent = message;
    element.style.display = 'block';
    setTimeout(() => {
        element.style.display = 'none';
    }, 5000);
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Auto-refresh logs every 5 seconds when on logs tab
let autoRefreshInterval;
function startAutoRefresh() {
    autoRefreshInterval = setInterval(() => {
        const logsTab = document.getElementById('logs-tab');
        if (logsTab && logsTab.style.display !== 'none') {
            loadLogs(currentPage);
        }
    }, 5000);
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    // Load initial configuration
    loadConfig();

    // Setup form submission
    document.getElementById('config-form').addEventListener('submit', saveConfig);

    // Set the initial view to the first configuration section
    const firstSubTab = document.querySelector('.sidebar .sub-tab');
    if (firstSubTab) {
        firstSubTab.click();
    }

    // Start auto-refresh for logs
    startAutoRefresh();
});
