document.addEventListener('DOMContentLoaded', function() {
    const chatContainer = document.getElementById('chatContainer');
    const chatInput = document.getElementById('chatInput');
    const sendButton = document.getElementById('sendMessage');
    const clearInputButton = document.getElementById('clearInput');
    const clearChatButton = document.getElementById('clearChat');
    const copyChatButton = document.getElementById('copyChat');
    const darkModeToggle = document.getElementById('darkModeToggle');
    const notificationBar = document.getElementById('notificationBar');

    const sessionId = Date.now()
    const backendHost = "backend.rag.svc.cluster.local"
    const backendPort = "8000"

    const socket = new WebSocket(`ws://${backendHost}:${backendPort}/ws/${sessionId}`);

    let chatContext = ""
    let thinkingMessage

    function showNotification(message) {
        notificationBar.textContent = message;
        notificationBar.style.display = 'block';
        setTimeout(() => {
            notificationBar.style.display = 'none';
        }, 3000);
    }

    function addMessage(content, isUser = true) {
        const messageGroup = document.createElement('div');
        messageGroup.classList.add('message-group', isUser ? 'user-group' : 'system-group');

        const label = document.createElement('div');
        label.classList.add('message-label');
        label.textContent = isUser ? 'User' : 'System';

        const messageElement = document.createElement('div');
        messageElement.classList.add('chat-message', isUser ? 'user-message' : 'system-message');

        const copyButton = document.createElement('button');
        copyButton.classList.add('copy-button');
        copyButton.textContent = '📋';
        copyButton.setAttribute('aria-label', 'Copy message');
        copyButton.addEventListener('click', function() {
            navigator.clipboard.writeText(content).then(function() {
                showNotification('Message copied to clipboard!');
            }, function(err) {
                console.error('Could not copy text: ', err);
            });
        });

        messageElement.appendChild(copyButton);

        // Check for code blocks and apply syntax highlighting
        const codeRegex = /```(\w+)?\s*([\s\S]*?)```/g;
        let lastIndex = 0;
        let match;

        while ((match = codeRegex.exec(content)) !== null) {
            // Add text before code block
            if (match.index > lastIndex) {
                messageElement.appendChild(document.createTextNode(content.slice(lastIndex, match.index)));
            }

            // Create code block
            const pre = document.createElement('pre');
            const code = document.createElement('code');
            if (match[1]) {
                code.className = `language-${match[1]}`;
            }
            code.textContent = match[2].trim();
            pre.appendChild(code);
            messageElement.appendChild(pre);

            lastIndex = match.index + match[0].length;
        }

        // Add remaining text after last code block
        if (lastIndex < content.length) {
            messageElement.appendChild(document.createTextNode(content.slice(lastIndex)));
        }

        messageGroup.appendChild(label);
        messageGroup.appendChild(messageElement);

        chatContainer.appendChild(messageGroup);
        chatContainer.scrollTop = chatContainer.scrollHeight;

        // Apply syntax highlighting
        messageElement.querySelectorAll('pre code').forEach((block) => {
            hljs.highlightBlock(block);
        });
    }

    function addThinkingMessage() {
        const messageGroup = document.createElement('div');
        messageGroup.classList.add('message-group', 'system-group');

        const label = document.createElement('div');
        label.classList.add('message-label');
        label.textContent = 'System';

        const messageElement = document.createElement('div');
        messageElement.classList.add('chat-message', 'system-message');

        const thinkingIndicator = document.createElement('div');
        thinkingIndicator.classList.add('thinking');

        const thinkingText = document.createTextNode('Thinking...');

        messageElement.appendChild(thinkingIndicator);
        messageElement.appendChild(thinkingText);

        messageGroup.appendChild(label);
        messageGroup.appendChild(messageElement);

        chatContainer.appendChild(messageGroup);
        chatContainer.scrollTop = chatContainer.scrollHeight;

        return messageGroup;
    }

    clearInputButton.addEventListener('click', function() {
        chatInput.value = '';
    });

    clearChatButton.addEventListener('click', function() {
        chatContainer.innerHTML = '';
        chatContext = "";
    });

    copyChatButton.addEventListener('click', function() {
        const chatText = Array.from(chatContainer.querySelectorAll('.message-group'))
            .map(group => {
                const label = group.querySelector('.message-label').textContent;
                const messageElement = group.querySelector('.chat-message');
                const message = Array.from(messageElement.childNodes)
                    .filter(node => node.nodeType === Node.TEXT_NODE || node.nodeName !== 'BUTTON')
                    .map(node => node.textContent.trim())
                    .join('');
                return `${label}: ${message}`;
            })
            .join('\n\n');
        navigator.clipboard.writeText(chatText).then(function() {
            showNotification('Chat copied to clipboard!');
        }, function(err) {
            console.error('Could not copy text: ', err);
        });
    });

    darkModeToggle.addEventListener('click', function() {
        document.body.classList.toggle('dark-mode');
        this.textContent = document.body.classList.contains('dark-mode') ? '☀️' : '🌙';
    });

    sendButton.addEventListener('click', function() {
        const message = chatInput.value.trim();
        if (message) {
            if (thinkingMessage != null ){
                chatContainer.removeChild(thinkingMessage);
            }

            addMessage(message, true);
            chatInput.value = '';

            chatContext += `${message} `;

            console.log("user:", message);
            console.log("chatContext", chatContext);

            // Query RAG
            fetch("http://localhost:8000/rag/query", {
                method: "POST",
                body: JSON.stringify({
                    sessionId: `${sessionId}`,
                    content: chatContext,
                }),
                headers: {
                    "Content-type": "application/json"
                }
            }).then(response => {response.json()
                .then(data => {console.log(data)}).catch(err => {console.error(err)})});

            thinkingMessage = addThinkingMessage();
        }
    });

    // Start WebSocket connection and display received messages
    socket.addEventListener("message", (event) => {
        const msg = JSON.parse(event.data);
        console.log("answer:", msg);
        if (msg) {
            if (thinkingMessage != null ){
                chatContainer.removeChild(thinkingMessage);
                thinkingMessage = null;
            }

            const chatMessage = document.createElement('div');
            chatMessage.classList.add('chat-message');
            chatInput.appendChild(chatMessage);
            chatInput.scrollTop = chatInput.scrollHeight;
            addMessage(msg.answer, false);
        }
    })

    chatInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendButton.click();
        }
    });

    // Add copy functionality to existing messages
    document.querySelectorAll('.copy-button').forEach(button => {
        button.addEventListener('click', function() {
            const messageText = this.nextSibling.textContent.trim();
            navigator.clipboard.writeText(messageText).then(function() {
                showNotification('Message copied to clipboard!');
            }, function(err) {
                console.error('Could not copy text: ', err);
            });
        });
    });
});
