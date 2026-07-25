type CatalogListProps<T extends { id: string }> = {
  items: T[];
  emptyMessage: string;
  renderItem: (item: T) => React.ReactNode;
};

export function CatalogList<T extends { id: string }>({
  items,
  emptyMessage,
  renderItem,
}: CatalogListProps<T>) {
  if (items.length === 0) {
    return <p className="text-sm text-muted-foreground">{emptyMessage}</p>;
  }

  return (
    <ul className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
      {items.map((item) => (
        <li key={item.id}>{renderItem(item)}</li>
      ))}
    </ul>
  );
}
