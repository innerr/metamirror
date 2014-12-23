package structs

// TODO: test Dict

/*
if __name__ == '__main__':
    d1 = Dictionary(111)
    d1.set(6, 9)
    d1.set(6, 3)

    s = StringIO()
    d1.commit(s.write)
    c = s.getvalue()
    d1.set(1, 2)
    d1.set(4, 8)
    d1.set(4, 0)
    s = StringIO()
    d1.commit(s.write)
    c = s.getvalue()
    
    d2 = Dictionary(222)
    d2.merge(c)
    assert set(d2.keys()) == set([1, 4])
    assert d2.get(1) == 2
    assert d2.get(4) == 0
    assert d2.history(4) == {111: 2}
    d2.set(4, 8)
    assert d2.get(4) == 8
    assert d2.history(4) == {111: 2, 222: 3}

    d2.set('child', Dictionary)
    s = d2.get('child')
    s.set('a', 'b')
    s.set('x', 'y')
    assert set(d2.keys()) == set([1, 4, 'child'])
    
    s = StringIO()
    d2.commit(s.write)
    c = s.getvalue()
    d3 = Dictionary(222)
    d3.merge(c)
    assert set(d3.keys()) == set([4, 'child'])
    assert d3.get(4) == 8
    s = d3.get('child')
    assert set(s.keys()) == set(['a', 'x'])
    assert s.get('a') == 'b'
    assert s.get('x') == 'y'

    d2.set('empty', Dictionary)
    assert set(d2.keys()) == set([1, 4, 'child', 'empty'])

    s = StringIO()
    d2.commit(s.write)
    c = s.getvalue()
    d3.merge(c)
    assert set(d3.keys()) == set([4, 'child', 'empty'])
    assert d3.get(4) == 8
    s = d3.get('child')
    assert set(s.keys()) == set(['a', 'x'])
    assert s.get('a') == 'b'
    assert s.get('x') == 'y'

    d3.merge(c)
    assert set(d3.keys()) == set([4, 'child', 'empty'])
    assert d3.get(4) == 8
    s = d3.get('child')
    assert set(s.keys()) == set(['a', 'x'])
    assert s.get('a') == 'b'
    assert s.get('x') == 'y'
*/
