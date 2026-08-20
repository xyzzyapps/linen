# Building a CSV Stream Summarizer

We start by defining the executable entry point at label `100`.

```python [app.py:100]
if __name__ == "__main__":
    main()
```

Before execution, we need system imports at label `10`. We can add inline snippets like `import sys` [app.py:10] directly in prose, followed by module-specific imports:

```python [app.py:12]
import csv
```

Next, we write our core calculation helper at label `30`.

```python [app.py:30]
def compute_avg(data):
    return sum(data) / len(data) if data else 0.0
```

Now we define `main()` at label `50`, which consumes `sys.stdin` and calls `compute_avg()`.

```python [app.py:50]
def main():
    reader = csv.reader(sys.stdin)
    vals = [float(row[0]) for row in reader if row]
    print(f"Average: {compute_avg(vals):.2f}")
```

We also define default configuration settings in a separate target file (`settings.json`):

```json [settings.json:10]
{
  "precision": 2,
  "stream_mode": true
}
```
