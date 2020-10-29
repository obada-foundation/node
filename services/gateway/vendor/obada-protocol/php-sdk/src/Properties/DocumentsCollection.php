<?php

declare(strict_types=1);

namespace Obada\Properties;

use IteratorAggregate;
use ArrayIterator;
use Obada\Properties\Document\Document;

class DocumentsCollection implements IteratorAggregate {

	use CollectionHash;

	protected array $items = [];

	public function __construct(Document ...$items) {
		$this->items = $items;
	}

	public function toArray() {
		return $this->items;
	}

	/**
	 * @param Document $document
	 * @return $this
	 */
	public function add(Document $document) {
		$this->items[] = $document;

		return $this;
	}

	public function getIterator() {
		return new ArrayIterator($this->items);
	}
}
