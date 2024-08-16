import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader
from torchvision import datasets, transforms

# Define the device to be used for training
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# Define your model, loss function, and optimizer
model = TransformerModel(emb_size=512, num_heads=8, num_layers=6, num_classes=10).to(device)
criterion = nn.CrossEntropyLoss()
optimizer = optim.Adam(model.parameters(), lr=0.001)

# Example dataset and dataloader setup
transform = transforms.Compose([transforms.ToTensor()])
train_dataset = datasets.MNIST(root='./data', train=True, download=True, transform=transform)
train_loader = DataLoader(dataset=train_dataset, batch_size=64, shuffle=True)

# Define the number of epochs
num_epochs = 10

# Training loop
for epoch in range(num_epochs):
    model.train()  # Set the model to training mode
    total_loss = 0

    for batch_idx, (inputs, targets) in enumerate(train_loader):
        inputs, targets = inputs.to(device), targets.to(device)

        # Forward pass
        outputs = model(inputs)
        loss = criterion(outputs, targets)

        # Backward pass and optimization
        optimizer.zero_grad()  # Clear gradients for this training step
        loss.backward()  # Propagation of the error back through the network
        optimizer.step()  # Update model parameters

        total_loss += loss.item()

        # Print statistics
        if (batch_idx + 1) % 100 == 0:
            print(f'Epoch [{epoch+1}/{num_epochs}], Step [{batch_idx+1}/{len(train_loader)}], Loss: {loss.item():.4f}')

    print(f'Epoch [{epoch+1}/{num_epochs}] finished with average loss: {total_loss / len(train_loader):.4f}')

# Optional: Save the trained model
torch.save(model.state_dict(), 'transformer_model.pth')
