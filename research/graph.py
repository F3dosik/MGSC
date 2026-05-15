import pandas as pd
import matplotlib.pyplot as plt

# 1. Загружаем данные из CSV (предполагается, что файл называется 'data.csv')
df = pd.read_csv('out.csv')

# 2. Создаем фигуру с двумя подграфиками (2 строки, 1 столбец)
fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(10, 8))

# 3. Первый график: nl от iter
ax1.plot(df['iter'], df['nl'], marker='o', linestyle='-', color='b', linewidth=2)
ax1.set_xlabel('Итерация')
ax1.set_ylabel('nl')
ax1.set_title('Зависимость nl от итерации')
ax1.grid(True, alpha=0.3)

# 4. Второй график: cost от iter
ax2.plot(df['iter'], df['cost'], marker='s', linestyle='-', color='r', linewidth=2)
ax2.set_xlabel('Итерация')
ax2.set_ylabel('cost')
ax2.set_title('Зависимость cost от итерации')
ax2.grid(True, alpha=0.3)

# 5. Настраиваем отображение
plt.tight_layout()
plt.show()