<?php

declare(strict_types=1);

namespace Obada\Properties;

use IteratorAggregate;
use ArrayIterator;
use Obada\Properties\StructuredData\Record;

class StructuredDataCollection implements IteratorAggregate {

	use CollectionHash;

	protected array $items = [];

	public function __construct(Record ...$items) {
		$this->items = $items;
	}

	public function toArray() {
		return $this->items;
	}

	/**
	 * @param Record $metadataRecord
	 * @return $this
	 */
	public function add(Record $metadataRecord) {
		$this->items[] = $metadataRecord;

		return $this;
	}

	public function getIterator() {
		return new ArrayIterator($this->items);
	}
}